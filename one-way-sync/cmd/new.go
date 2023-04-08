package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ahui2016/go-scripts/util"
)

const tomlName = "one-way-sync.toml"

var overwrite bool

type Pair struct {
	Src string
	Dst string
}

func (p *Pair) String() string {
	src := "SrcDir: " + p.Src
	dst := "DstDir: " + p.Dst
	return fmt.Sprintf("%s\n%s", src, dst)
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new SrcDir DstDir",
	Short: "Create a new config file (one-way-sync.toml).",
	Example: `
    one-way-sync new '.' 'D:/temp'`,
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := util.CheckOverwriteFile(tomlName, overwrite); err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		srcDir, e1 := filepath.Abs(args[0])
		dstDir, e2 := filepath.Abs(args[1])
		if err := util.WrapErrors(e1, e2); err != nil {
			log.Fatal(err)
		}
		pair := Pair{srcDir, dstDir}
		lo.Must0(util.WriteTOML(pair, tomlName))
		fmt.Println("Write " + tomlName)
		fmt.Println(pair.String())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().BoolVar(
		&overwrite,
		"overwrite",
		false,
		"Overwrite disabled by default, set to true to overwrite files",
	)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
