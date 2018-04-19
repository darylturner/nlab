package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var jsonOut bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nlab",
	Short: "A tool for making network labs under KVM simpler",
	Long: `nlab can be used to create Linux bridges and taps to
simulate complicated network topologies and launch KVM
virtual-machines with sane defaults.`,
	Version: "0.8.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)
		if jsonOut == true {
			log.SetFormatter(&log.JSONFormatter{})
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "Output formatted as JSON to stdout")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
