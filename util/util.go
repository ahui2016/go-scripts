package util

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
)

const (
	ReadonlyFilePerm = 0555
	NormalFilePerm   = 0666
	NormalFolerPerm  = 0750
)

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

func PathIsNotExist(name string) (ok bool) {
	_, err := os.Lstat(name)
	if os.IsNotExist(err) {
		ok = true
		err = nil
	}
	lo.Must0(err)
	return
}

func PathIsExist(name string) bool {
	return !PathIsNotExist(name)
}

// CheckOverwriteFile Returns nil if overwrite is true,
// or an error message if overwrite is false and the file already exists
// (file overwriting is forbidden).
func CheckOverwriteFile(name string, overwrite bool) error {
	if !overwrite && PathIsExist(name) {
		return fmt.Errorf("file exists: %s", name)
	}
	return nil
}
