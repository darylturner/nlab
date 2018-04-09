package cmd

import (
	"fmt"
	"os/exec"

	"github.com/darylturner/nlab/config"
	"github.com/darylturner/nlab/network"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <config.yml>",
	Short: "Create Linux bridge/taps",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.ParseConfig(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"config": args[0],
				"error":  err,
			}).Fatal("error parsing configuration")
		}

		netMap, err := network.GetMap(cfg)
		if err != nil {
			log.Fatal(err)
		}

		for _, net := range netMap {
			if err := net.Create(); err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"net": net,
				}).Error("error creating vm network")
			}
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	networkCmd.AddCommand(createCmd)
}
