// Package cmd 命令入口
package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:"mars",
	Short:"HTTP(S) proxy",
}

// Execute 执行命令
func Execute() {
	err := rootCmd.Execute()
	if err != nil {

	}
}
