package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ahui2016/go-scripts/util"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new SrcDir DstDir",
	Short: "Create a new config file (toml).",
	Example: `
    one-way-sync new '.' 'D:/temp'`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		srcDir, e1 := filepath.Abs(args[0])
		dstDir, e2 := filepath.Abs(args[1])
		if err := util.WrapErrors(e1, e2); err != nil {
			return err
		}
		fmt.Println(srcDir)
		fmt.Println(dstDir)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
