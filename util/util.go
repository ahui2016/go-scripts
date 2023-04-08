package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"golang.org/x/crypto/blake2b"
)

const (
	ReadonlyFilePerm = 0555
	NormalFilePerm   = 0666
	NormalFolerPerm  = 0750
)

type HexString = string

// WrapErrors 把多个错误合并为一个错误.
func WrapErrors(allErrors ...error) (wrapped error) {
	for _, err := range allErrors {
		if err != nil {
			if wrapped == nil {
				wrapped = err
			} else {
				wrapped = fmt.Errorf("%w | %w", wrapped, err)
			}
		}
	}
	return
}

// WriteFile 写檔案, 如果 perm 等于零, 则使用默认权限.
func WriteFile(name string, data []byte, perm fs.FileMode) error {
	if perm == 0 {
		perm = NormalFilePerm
	}
	return os.WriteFile(name, data, perm)
}

func WriteTOML(data interface{}, filename string) error {
	dataTOML, err := toml.Marshal(data)
	if err != nil {
		return err
	}
	return WriteFile(filename, dataTOML, 0)
}

// WriteJSON 把 data 转换为漂亮格式的 JSON 并写入檔案 filename 中。
func WriteJSON(data interface{}, filename string) error {
	dataJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return WriteFile(filename, dataJSON, 0)
}

func PathIsNotExist(name string) bool {
	_, ok := PathIsExist(name)
	return !ok
}

func PathIsExist(name string) (fs.FileInfo, bool) {
	info, err := os.Lstat(name)
	if os.IsNotExist(err) {
		return nil, false
	}
	lo.Must0(err)
	return info, true
}

// CheckOverwriteFile Returns nil if overwrite is true,
// or an error message if overwrite is false and the file already exists
// (file overwriting is forbidden).
func CheckOverwriteFile(name string, overwrite bool) error {
	filePath, err := filepath.Abs(name)
	if err != nil {
		return err
	}
	_, exist := PathIsExist(name)
	if !overwrite && exist {
		return fmt.Errorf("file exists: %s", filePath)
	}
	return nil
}

// BLAKE2b is faster than MD5, SHA-1, SHA-2, and SHA-3, on 64-bit x86-64 and ARM architectures.
// https://en.wikipedia.org/wiki/BLAKE_(hash_function)#BLAKE2
// https://blog.min.io/fast-hashing-in-golang-using-blake2/
// https://pkg.go.dev/crypto/sha256#example-New-File
func FileSum512(name string) (HexString, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := lo.Must(blake2b.New512(nil))
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	checksum := h.Sum(nil)
	return hex.EncodeToString(checksum), nil
}

// https://stackoverflow.com/questions/30376921/how-do-you-copy-a-file-in-go
func CopyFile(dstPath, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err1 := io.Copy(dst, src)
	err2 := dst.Sync()
	return WrapErrors(err1, err2)
}
