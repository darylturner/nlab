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
		var conf string
		if len(args) > 0 {
			conf = args[0]
		} else {
			conf = "-"
		}
		cfg, err := config.Parse(conf)
		if err != nil {
			log.WithFields(log.Fields{
				"config": conf,
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
	Args: cobra.MaximumNArgs(1),
}

func init() {
	networkCmd.AddCommand(destroyCmd)
}
