package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Used for flags.
var ()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "one-way-sync",
	Short: "單向同步資料夾",
	Long: `one-way-sync (單向同步資料夾):

只同步第一層檔案, 不同步子目錄.
以 SrcDir 為準, 向 DstDir 添加檔案, 或更新/刪除 DstDir 中的檔案.
使用方法: 'one-way-sync config.toml', 其中 config.toml 的內容如下.

    SrcDir: '/path/to/src-dir'
    DstDir: '/path/to/dst-dir'
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
