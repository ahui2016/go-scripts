package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ahui2016/go-scripts/util"
)

var Separator = string(filepath.Separator)

var config = new(Config)

var forceRun bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "one-way-sync",
	Short: "單向同步資料夾",
	Long: `one-way-sync (單向同步資料夾):

只同步第一層檔案, 不同步子目錄.
以 Src 為準, 向 Dst 添加檔案, 或更新/刪除 Dst 中的檔案.
使用方法舉例: 'one-way-sync config.toml', 其中 config.toml 的內容如下.

    Src = '/path/to/src-dir'
    Dst = '/path/to/dst-dir'
`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		config.Load(args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !forceRun {
			fmt.Printf("\n現在是 **Dry Run** 模式, 僅列印將要發生的變化.\n")
			fmt.Print("只有使用 --force 參數才會實際執行.\n")
		}
		fmt.Printf("\n[Source] %s\n", config.Src)
		fmt.Printf("[Target] %s\n\n", config.Dst)
		n, err := addOrUpdateFiles(config.Src, config.Dst, forceRun)
		if err != nil {
			log.Fatal(err)
		}
		m, err := deleteFiles(config.Src, config.Dst, forceRun)
		if err != nil {
			log.Fatal(err)
		}
		if n+m == 0 {
			fmt.Println("兩個資料夾內的檔案相同 (本程式不處理子目錄)")
		}
		fmt.Println()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.one-way-sync.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVar(
		&forceRun,
		"force",
		false,
		"默認僅列印將要發生變化的檔案, 只有使用該參數才會實際執行.",
	)

}

func addOrUpdateFiles(srcDir, dstDir string, force bool) (count int, err error) {
	files, err := filepath.Glob(srcDir + Separator + "*")
	if err != nil {
		return
	}
	for _, srcFile := range files {
		info, err := os.Lstat(srcFile)
		if err != nil {
			return count, err
		}

		// 跳过资料夹
		if info.IsDir() {
			continue
		}

		dstFile := filepath.Join(dstDir, info.Name())
		_, err = os.Lstat(dstFile)
		dstNotExist := os.IsNotExist(err)
		if dstNotExist {
			err = nil
		}
		if err != nil {
			return count, err
		}

		// 新增文档
		if dstNotExist {
			if force {
				if err := util.CopyFile(dstFile, srcFile); err != nil {
					return count, err
				}
			}
			fmt.Printf("ADD => %s\n", dstFile)
			count++
			continue
		}

		// 对比文档, 覆盖文档
		srcSum, e1 := util.FileSum512(srcFile)
		dstSum, e2 := util.FileSum512(dstFile)
		if err := util.WrapErrors(e1, e2); err != nil {
			return count, err
		}
		if srcSum != dstSum {
			if force {
				if err := util.CopyFile(dstFile, srcFile); err != nil {
					return count, err
				}
			}
			fmt.Printf("UPDATE => %s\n", dstFile)
			count++
		}
	}
	return
}

func deleteFiles(srcDir, dstDir string, force bool) (count int, err error) {
	files, err := filepath.Glob(dstDir + Separator + "*")
	if err != nil {
		return
	}
	for _, dstFile := range files {
		info, err := os.Lstat(dstFile)
		if err != nil {
			return count, err
		}

		// 跳过资料夹
		if info.IsDir() {
			continue
		}

		srcFile := filepath.Join(srcDir, info.Name())
		_, err = os.Lstat(srcFile)
		srcNotExist := os.IsNotExist(err)
		if srcNotExist {
			err = nil
		}
		if err != nil {
			return count, err
		}

		if srcNotExist {
			if force {
				if err := os.Remove(dstFile); err != nil {
					return count, err
				}
			}
			fmt.Printf("DELETE => %s\n", dstFile)
			count++
		}
	}
	return
}
