package cmd

import (
	"github.com/ouqiang/mars/internal/common/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:"version",
	Short:"print version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(version.Format())
	},
}

func init()  {
	rootCmd.AddCommand(versionCmd)
}