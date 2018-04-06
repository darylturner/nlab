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
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <config.yml>",
	Short: "Run virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := parseConfig(args[0])
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

		serialPortBase := 50000 // need to make this more dynamic
		for nodeIdx, node := range cfg.Nodes {
			specifiedTag := cmd.Flag("tag").Value.String()
			if specifiedTag == "" || specifiedTag == node.Tag {
				telnetPort := serialPortBase + nodeIdx

				qemuArgs := []string{
					"-name", node.Tag, "-daemonize",
					"-smp", node.Resources.CPU,
					"-pidfile", fmt.Sprintf("/var/run/nlab/%v/%v.pid", cfg.Tag, node.Tag),
					"-m", node.Resources.Memory,
					"-display", "none", "-serial", fmt.Sprintf("telnet::%v,nowait,server", telnetPort),
				}

				for _, disk := range node.Resources.Disks {
					cmd := []string{"-drive", fmt.Sprintf("format=%v,file=%v", disk.Format, disk.File)}
					qemuArgs = append(qemuArgs, cmd...)
				}

				if node.Resources.CDROM != "" {
					qemuArgs = append(qemuArgs, "-cdrom "+node.Resources.CDROM)
				}

				if node.Network.Management == true {
					tapName := fmt.Sprintf("mng%s", node.Tag)
					qemuArgs = append(qemuArgs, linkCmd(cfg.ManagementBridge, tapName)...)
				}
				for _, link := range node.Network.Links {
					tapName := fixedLengthTap(link, node.Tag)
					qemuArgs = append(qemuArgs, linkCmd(link, tapName)...)
				}

				dryRun := cmd.Flag("no-launch").Value.String()
				if strings.ToLower(dryRun) != "true" {
					if err := exec.Command("kvm", qemuArgs...).Run(); err != nil {
						log.WithFields(log.Fields{
							"tag":   node.Tag,
							"error": err,
						}).Error("error starting node")
					} else {
						log.WithFields(log.Fields{
							"tag":    node.Tag,
							"serial": fmt.Sprintf("telnet://localhost:%v", telnetPort),
						}).Info("running")
					}
				} else {
					fmt.Println("kvm " + strings.Join(qemuArgs, " "))
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

func linkCmd(link, tap string) []string {
	return []string{
		"-device", fmt.Sprintf("e1000,netdev=%s,mac=%s", link, generateMAC()),
		"-netdev", fmt.Sprintf("tap,id=%s,ifname=%s,script=no", link, tap),
	}
}

func randomByte() int {
	return rand.Intn(256)
}

func generateMAC() string {
	baseMAC := "52:54:00"
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(
		baseMAC+":%02x:%02x:%02x",
		randomByte(),
		randomByte(),
		randomByte(),
	)
}
