package cmd

import (
	"fmt"

	"github.com/siliconandsolder/gorp/grep"
	"github.com/spf13/cobra"
)

var verboseMode bool = false

var RootCmd = &cobra.Command{
	Use:   "gorp",
	Short: "Greppin' with Go!",
	Run: func(cmd *cobra.Command, args []string) {
		gm := grep.NewGrepManager(verboseMode, args[1], args[2], args[3:])
		fmt.Printf("%v", gm)
	},
}

func init() {
	RootCmd.Flags().BoolVarP(&verboseMode, "verbose", "v", false, "run gorp in verbose mode")
}
