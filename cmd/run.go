package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/darylturner/nlab/config"
	"github.com/darylturner/nlab/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <config.yml>",
	Short: "Run virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Parse(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"config": args[0],
				"error":  err,
			}).Fatal("error parsing configuration")
		}

		if err := os.MkdirAll(fmt.Sprintf("/var/run/nlab/%v", cfg.Tag), os.ModePerm); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("unable to create pid dir")
		}

		tag := cmd.Flag("tag").Value.String()
		noLaunch := cmd.Flag("no-launch").Value.String()
		dryRun, err := strconv.ParseBool(noLaunch)
		if err != nil {
			panic(err)
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

				if err := nd.Run(cfg, dryRun); err != nil {
					log.WithFields(log.Fields{
						"tag": ndConf.Tag,
						"err": err,
					}).Error("node run failed")
					continue
				}
			}
		}

	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("no-launch", "n", false, "Don't launch machines. Output qemu command to stdout.")
	runCmd.Flags().StringP("tag", "t", "", "Launch only virtual machine matching tag.")
}
