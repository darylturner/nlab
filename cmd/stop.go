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
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <config.yml>",
	Short: "Stop virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := parseConfig(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"config": args[0],
				"error":  err,
			}).Fatal("error parsing configuration")
		}

		pidPath := fmt.Sprintf("/var/run/nlab/%v/", cfg.Tag)

		specifiedTag := cmd.Flag("tag").Value.String()
		for _, node := range cfg.Nodes {
			if specifiedTag == "" || specifiedTag == node.Tag {
				pidBytes, err := ioutil.ReadFile(pidPath + fmt.Sprintf("%v.pid", node.Tag))
				if err != nil {
					log.WithFields(log.Fields{
						"error":    err,
						"tag":      node.Tag,
						"pid_path": pidPath,
					}).Error("unable to read pid")
					continue
				}

				pidString := strings.TrimSpace(string(pidBytes))
				pid, err := strconv.Atoi(pidString)
				if err != nil {
					panic(err)
				}

				proc, err := os.FindProcess(pid)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"tag":   node.Tag,
					}).Error("unable to find process")
					continue
				}

				if err := proc.Kill(); err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"tag":   node.Tag,
					}).Error("unable to kill process")
				} else {
					log.WithFields(log.Fields{
						"tag": node.Tag,
					}).Info("node stopped")
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
