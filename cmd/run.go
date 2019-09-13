package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/darylturner/nlab/internal/config"
	"github.com/darylturner/nlab/internal/network"
	"github.com/darylturner/nlab/internal/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <config.json>",
	Short: "Run virtual machines",
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

		if err := os.MkdirAll(fmt.Sprintf("/var/run/nlab/%v", cfg.Tag), os.ModePerm); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("unable to create pid dir")
		}

		tag := cmd.Flag("tag").Value.String()
		nlFlag := cmd.Flag("no-launch").Value.String()

		dryRun, err := strconv.ParseBool(nlFlag)
		if err != nil {
			panic(err)
		}

		pwMap := make(map[string]*network.PseudoWire)
		if cfg.PseudoWire {
			pwMap, err = network.GetPseudoWireMap(cfg)
			if err != nil {
				log.Fatal("error creating pseudowire map")
			}
		}

		for _, ndConf := range cfg.Nodes {
			if tag == ndConf.Tag || tag == "" {
				nd, err := node.New(ndConf)
				if err != nil {
					log.WithFields(log.Fields{
						"tag": ndConf.Tag,
						"err": err,
					}).Error("error creating node object")
					continue
				}

				status, err := nd.Run(cfg, dryRun, pwMap)
				if err != nil {
					log.WithFields(log.Fields{
						"tag": ndConf.Tag,
						"err": err,
					}).Error("node run failed")
					continue
				}

				if status != nil {
					log.WithFields(status).Info("running")
				} else {
					if !dryRun {
						log.WithFields(log.Fields{"tag": ndConf.Tag}).Info("running")
					}
				}
			}
		}

	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("no-launch", "n", false, "Don't launch machines. Output qemu command to stdout.")
	runCmd.Flags().StringP("tag", "t", "", "Launch only virtual machine matching tag.")
}
