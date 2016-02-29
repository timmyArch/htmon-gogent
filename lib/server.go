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
  "net/http"
  "net/url"
  "os"
  "strconv"
  "io/ioutil"
  "encoding/json"

  log "github.com/Sirupsen/logrus"
)

var SpoofedHostname string

type Server struct {
  User string
  Password string
  ApiUrl string
}

func (s *Server) PushCheck(c *Check) error {
  values := url.Values{}
  values.Add("metric", c.Metric)
  values.Add("expire", strconv.Itoa(c.Expire))
  values.Add("hostname", s.Hostname())
  values.Add("value", c.RunCommand())

  _, err := s.Request("/metrics", values)

  return err
}

func (s *Server) Keepalive() error {
  values := url.Values{}
  values.Add("expire", "10")
  values.Add("hostname", s.Hostname())

  _, err := s.Request("/metrics/keepalive", values)

  return err
}

func (s *Server) Schema() map[string]interface{} {
  values := url.Values{}
  values.Add("hostname", s.Hostname())

  request, _ := http.NewRequest("GET", s.ApiUrl+"/metrics/schema?"+values.Encode(), nil)
  request.SetBasicAuth(s.User, s.Password)
  client := &http.Client{}
  resp, _ := client.Do(request)

  var schema map[string]interface{}

  defer resp.Body.Close()

  body, _ := ioutil.ReadAll(resp.Body)
  
  e := json.Unmarshal(body, &schema)

  log.WithFields(log.Fields{
    "Schema": schema,
    "JSONError": e,
  }).Debug("Loaded schema.")

  return schema
}

func (s *Server) Hostname() string {
  if SpoofedHostname != "" {
    return SpoofedHostname
  } else {
    hostname, _ := os.Hostname()
    return hostname
  }
}

func (s *Server) Request(path string, values url.Values) (*http.Response, error) {
  request, err := http.NewRequest("POST", s.ApiUrl+path+"?"+values.Encode(), nil)
  
  log.WithFields(log.Fields{
      "Url": request.URL,
  }).Debug("Starting HTTP API call.")

  if err != nil {
    log.WithFields(log.Fields{
        "Request": request,
        "Error": err,
    }).Error("Retrieving new request failed.")
  }

  request.SetBasicAuth(s.User, s.Password)
  client := &http.Client{}
  resp, err := client.Do(request)

  if err != nil {
    log.WithFields(log.Fields{
        "Request": request,
        "Response": resp,
        "Error": err,
    }).Error("Launching http request failed.")
  }
  
  log.WithFields(log.Fields{
      "Response": resp.Status,
  }).Debug("Api call done.")


  return resp, err
}
