// Copyright 2016 Tim Foerster <github@mailserver.1n3t.de>
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
	"os"

	"github.com/spf13/cobra"
)

var ConfigFile string
var Version bool

var RootCmd = &cobra.Command{
	Use:   "htmon-gogent",
	Short: "Htmon client side agent written in Go",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
    if Version {
      fmt.Println("Version: 0.1.0")
      os.Exit(0)
    }
  },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "", 
    "config file (default is $HOME/.htmon-gogent.yaml)")
	RootCmd.Flags().BoolVarP(&Version, "version", "v", false, "Print version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
