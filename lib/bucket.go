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

package lib

import (
  log "github.com/Sirupsen/logrus"
)

type Bucket struct {
  Server *Server
  Checks []*Check
}

func (b *Bucket) Process() {
  var processed_checks []CheckResult

  checks := b.SelectChecks()

  channel := make(chan CheckResult)

  b.Server.Keepalive()
  log.Info(len(b.Checks), " bucket available checks ...")
  log.Info(len(checks), " preselected checks ...")

  for i, check := range checks {
    log.WithFields(log.Fields{
      "CheckMetric": check.Metric,
      "Run": i + 1,
    }).Debug("Start routine ...")
    go func(c chan<- CheckResult, check *Check) {
      c <- check.Process() 
    }(channel, check)
  }

  for i := 0; i < len(checks); i++ {
    tmp_check := <-channel
    processed_checks = append(processed_checks, tmp_check)
    log.WithFields(log.Fields{
      "Run": i+1,
      "Amount": len(checks),
      "CheckMetric": tmp_check.Metric,
    }).Debug("Check processed...")
  }

  log.WithFields(log.Fields{
    "ProcessedChecks": len(processed_checks),
  }).Debug("Checks processed ...")
}

func (b *Bucket) SelectChecks() []*Check {
  var preselected_checks []*Check

  for _, check := range b.Checks {
    if check.ReadyForProcess() { 
      preselected_checks = append(preselected_checks, check)
    }
  }

  return preselected_checks
}
