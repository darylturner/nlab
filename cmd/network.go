package cmd

import (
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Create and destroy Linux bridge/taps",
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
