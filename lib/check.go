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
  "time"
  "strings"
  "os/exec"

  "gopkg.in/yaml.v2"
  log "github.com/Sirupsen/logrus"
)

/*
Check object definition
*/

type Check struct {
  Server *Server
  Metric string
  Interval int64
  Expire int
  Command string
  LastRun int64
}

type UnparsedCheck struct {
  Check
  Placeholder interface{}
}

type CheckResult struct {
  Check
  Error error
}

func (o *Check) Process() CheckResult {
  err := o.Server.PushCheck(o)
  o.LastRun = time.Now().Unix()
  return CheckResult{Check: *o, Error: err,}
}

func (o *Check) ReadyForProcess() bool {
  t := time.Now().Unix()
  return o.LastRun + o.Interval < t
}

func (o *Check) RunCommand() string {
  cmd, err := exec.Command("/bin/sh", "-c", o.Command).Output()

  tcmd := strings.Trim(string(cmd), "\n")

  if err != nil {
    log.WithFields(log.Fields{
      "Command": o.Command,
      "Return": tcmd,
      "Error": err,
    }).Error("Command failed.")
  }

  return tcmd
}

func (o *UnparsedCheck) AnyPlaceholder() (bool, []string)  {
  orig, ok := o.Placeholder.([]interface{})
  var str []string
  if ok {
    for _, item := range orig {
      str = append(str, item.(string))
    }
  } else {
    placeholders, ok := o.Placeholder.([]string)
    if ok {
      str = placeholders
    }
  }
  return (len(str) > 0), str
}

func (o *UnparsedCheck) ParsedMetric(p string) string {
  return HandleExtractedObject(o.Metric, p)
}

func (o *UnparsedCheck) ParsedCommand(p string) string {
  return HandleExtractedObject(o.Command, p)
}

func HandleExtractedObject(c string, p string) string {
  return strings.Replace(c, "$$placeholder$$", p, -1)
}


/*
Check collection object definition
*/

type CheckCollection struct {
  Server Server
  Checks []UnparsedCheck
}

func (f *CheckCollection) UnmarshalYAML(bs []byte) error {
  return yaml.Unmarshal(bs, &f.Checks)
}

func (f *CheckCollection) PreProcess() []*Check {

  var pre_processed []*Check
  var placeholder_count, schema_count, single_count int

  log.Debug("Start pre-processing check collection.")

  for _, check := range f.Checks {
    ok, placeholders := check.AnyPlaceholder()
    if ok {
      for _, placeholder := range placeholders {
        placeholder_count++
        new_check := Check{
          Server: &f.Server,
          Interval: check.Interval,
          Expire: check.Expire,
        }
        new_check.Metric = check.ParsedMetric(placeholder)
        new_check.Command = check.ParsedCommand(placeholder)
        pre_processed = append(pre_processed, &new_check)
      }
    } else {
      c := CheckCollection{
        Server: f.Server,
      }
      if check.Placeholder == nil {
        single_count++
        new_check := Check{
          Server: &f.Server,
          Interval: check.Interval,
          Expire: check.Expire,
          Metric: check.Metric,
          Command: check.Command,
        }
        pre_processed = append(pre_processed, &new_check)
      } else {
        orig, ok := check.Placeholder.(string)
        schema := f.Server.Schema()
        if ok && schema[orig] != nil {
          var str []string
          enforced_array, parse_ok := schema[orig].([]interface{})
          if parse_ok {
            for _, item := range enforced_array {
              schema_count++
              str = append(str, item.(string))
            }
            check.Placeholder = str
            c.Checks = []UnparsedCheck{check}
            for _, t := range c.PreProcess() {
              pre_processed = append(pre_processed, t)
            }
          }
        }
      }
    }
  }
  log.WithFields(log.Fields{
    "AmountOfChecks": len(pre_processed),
    "ChecksBySchema": schema_count,
    "ChecksByPlaceholder": placeholder_count-schema_count,
    "ChecksByOther": single_count,
  }).Debug("Checks successfully pre-processed.")

  return pre_processed
}
