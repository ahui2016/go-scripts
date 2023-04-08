package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ahui2016/go-scripts/util"
)

const configFileName = "one-way-sync.toml"

var overwrite bool

type Config struct {
	Src string
	Dst string
}

func (cfg *Config) String() string {
	src := "SrcDir: " + cfg.Src
	dst := "DstDir: " + cfg.Dst
	return fmt.Sprintf("%s\n%s", src, dst)
}

func (cfg *Config) Load(tomlPath string) {
	if tomlPath == "" {
		tomlPath = configFileName
	}
	data := lo.Must(os.ReadFile(tomlPath))
	lo.Must0(toml.Unmarshal(data, cfg))
	cfg.check()
}

func (cfg *Config) check() {
	for _, dirPath := range []string{cfg.Src, cfg.Dst} {
		info, ok := util.PathIsExist(dirPath)
		if !ok {
			log.Fatalf("Not Found: %s\n", dirPath)
		}
		if !info.IsDir() {
			log.Fatalf("不是資料夾: %s\n", dirPath)
		}
		if !filepath.IsAbs(dirPath) {
			log.Fatalf("不是絕對路徑: %s\n", dirPath)
		}
	}
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new SrcDir DstDir",
	Short: "Create a new config file (one-way-sync.toml).",
	Example: `
    one-way-sync new '.' 'D:/temp'`,
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := util.CheckOverwriteFile(configFileName, overwrite); err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		srcDir, e1 := filepath.Abs(args[0])
		dstDir, e2 := filepath.Abs(args[1])
		if err := util.WrapErrors(e1, e2); err != nil {
			log.Fatal(err)
		}
		cfg := Config{srcDir, dstDir}
		lo.Must0(util.WriteTOML(cfg, configFileName))
		fmt.Println("Write " + configFileName)
		fmt.Println(cfg.String())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().BoolVar(
		&overwrite,
		"overwrite",
		false,
		"默認禁止覆蓋檔案, 使用該參數則允許覆蓋.",
	)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
