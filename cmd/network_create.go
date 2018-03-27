// Copyright © 2018 Daryl Turner <daryl@layer-eight.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <config.yml>",
	Short: "Create Linux bridge/taps",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := parseConfig(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"config": args[0],
				"error":  err,
			}).Fatal("error parsing configuration")
		}

		// create bridges for each segment
		for _, br := range cfg.Links {
			if err := createBridge(br.Tag); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"bridge": br.Tag,
				}).Error("error creating bridge")
			}
		}

		// create taps and add to bridges
		for _, node := range cfg.Nodes {
			if node.Network.Management == true {
				mngBridge := cfg.ManagementBridge
				tapName := fmt.Sprintf("mng%s", node.Tag)
				if err := createTap(tapName, mngBridge); err != nil {
					log.WithFields(log.Fields{
						"error":  err,
						"bridge": mngBridge,
						"tap":    tapName,
						"node":   node.Tag,
					}).Error("error creating management tap")
				}
			}

			for _, link := range node.Network.Links {
				segmentBr := link
				tapName := fixedLengthTap(link, node.Tag)
				if err := createTap(tapName, segmentBr); err != nil {
					log.WithFields(log.Fields{
						"error":  err,
						"bridge": segmentBr,
						"tap":    tapName,
						"node":   node.Tag,
					}).Error("error creating segment tap")
				}
			}
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	networkCmd.AddCommand(createCmd)
}

func createTap(name, bridge string) error {
	contextLog := log.WithFields(log.Fields{
		"tap": name, "bridge": bridge,
	})

	contextLog.Info("creating tap")
	if err := exec.Command("ip", "tuntap", "add", "dev", name, "mode", "tap").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "set", name, "master", bridge).Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "set", name, "up").Run(); err != nil {
		return err
	}

	return nil
}

func createBridge(name string) error {
	contextLog := log.WithFields(log.Fields{
		"bridge": name,
	})

	contextLog.Info("creating bridge")
	if err := exec.Command("ip", "link", "add", name, "type", "bridge").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "set", name, "up").Run(); err != nil {
		return err
	}

	return nil
}
