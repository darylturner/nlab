package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nlab",
	Short: "A tool for making network labs under KVM simpler",
	Long: `nlab can be used to create Linux bridges and taps to
simulate complicated network topologies and launch KVM
virtual-machines with sane defaults.`,
	Version: "0.5.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
