package vmtool

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vmtool",
	Short: "VMTool is a terminal-based VM streaming tool",
	Long:  `A terminal-based alternative to UTM that provisions VMs and streams them to a browser.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
}
