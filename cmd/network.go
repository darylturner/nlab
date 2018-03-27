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
	"crypto/sha256"
	"fmt"

	"github.com/spf13/cobra"
)

func fixedLengthTap(link, node string) string {
	// required as linux has a character limit on tap names,
	// we need to be fairly sure of deriving a unique tap name
	// to avoid name collisions.
	hash := sha256.New()
	hash.Write([]byte(link))
	hash.Write([]byte(node))

	return fmt.Sprintf("%.5x_tap0", hash.Sum(nil))
}

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Create and destroy Linux bridge/taps",
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
