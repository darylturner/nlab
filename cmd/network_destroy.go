package cmd

import (
	"github.com/darylturner/nlab/internal/config"
	"github.com/darylturner/nlab/internal/network"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy <config.json>",
	Short: "Destroy Linux bridge/taps",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Parse(args[0])
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
			if err := net.Destroy(); err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"net": net,
				}).Error("error destroying vm network")
			}
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	networkCmd.AddCommand(destroyCmd)
}
