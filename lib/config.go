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
  "io/ioutil"
  "gopkg.in/yaml.v2"

  log "github.com/Sirupsen/logrus"
)

type Config struct {
  Path string
  LoadedCollection *CheckCollection
}

func (c *Config) LoadFile() []byte {
  log.WithFields(log.Fields{
      "Path": c.Path,
  }).Debug("Reading config ...")

  data, e := ioutil.ReadFile(c.Path)
  c.Error(e)

  log.Debug("Config successfully read.")

  return data
}

func (c *Config) ParseConfig() {
  if c.LoadedCollection == nil {
    checks := CheckCollection{}

    log.Debug("Start parsing check collection.")

    e := yaml.Unmarshal(c.LoadFile(), &checks)
    c.Error(e)

    c.LoadedCollection = &checks

    log.WithFields(log.Fields{
        "CheckCollection": checks,
    }).Debug("Checks successfully parsed.")
  }
}

func (c *Config) Checks() []*Check {
  c.ParseConfig()
  return c.LoadedCollection.PreProcess()
}

func (c *Config) Server() *Server {
  c.ParseConfig()
  return &c.LoadedCollection.Server
}

func (c *Config) Error(e error) {
  if e != nil {
    log.Panic(e)
  }
}
