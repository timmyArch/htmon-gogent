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
  "time"
  "encoding/json"
  "fmt"
  "os"

	"github.com/spf13/cobra"
  "github.com/timmyArch/htmon-gogent/lib"
  log "github.com/Sirupsen/logrus"
)

var TestRun bool
var SchemaOnly bool

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "agent",
	Short: "Handle agent actions",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
    log.SetOutput(os.Stderr)

    c := lib.Config{
      Path: ConfigFile,
    }
    b := lib.Bucket{
      Checks: c.Checks(),
      Server: c.Server(),
    }
    schema := b.Server.Schema()
    if TestRun {
      log.SetLevel(log.DebugLevel)
      b.Process()
    } else if SchemaOnly {
      log.SetLevel(log.PanicLevel)
      moo, _ := json.MarshalIndent(schema, "", "  ")
      fmt.Println(string(moo))
    } else {
      log.SetLevel(log.DebugLevel)
      for ;; {
        b.Process()
        time.Sleep(time.Second * 1)
      }
    }
	},
}

func init() {
	RootCmd.AddCommand(launchCmd)
	launchCmd.Flags().BoolVarP(&TestRun, "test", "t", false, "Testrun, just one report.")
	launchCmd.Flags().BoolVarP(&SchemaOnly, "schema", "s", false, "Show loaded schema.")
	launchCmd.Flags().StringVarP(&lib.SpoofedHostname, "spoof-hostname", "", "", 
    "Spoof given hostname instead of current one.")
}
