package cmd

import (
	"github.com/darylturner/nlab/config"
	"github.com/darylturner/nlab/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <config.yml>",
	Short: "Stop virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Parse(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"config": args[0],
				"error":  err,
			}).Fatal("error parsing configuration")
		}

		tag := cmd.Flag("tag").Value.String()
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

				if err := nd.Stop(cfg); err != nil {
					log.WithFields(log.Fields{
						"tag": ndConf.Tag,
						"err": err,
					}).Error("node stop failed")
					continue
				}
			}
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringP("tag", "t", "", "Stop only virtual machine matching tag.")
}
