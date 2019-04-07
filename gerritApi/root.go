// Copyright © 2019 Robert Sotomski <sotomski@gmail.com>
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

package gerritApi

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/utils"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"time"
)

type Project struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func Initialize() {
	resty.SetHostURL(viper.GetString("GOOPSC_GERRIT_URL"))
	resty.SetTimeout(1 * time.Minute)

	// Headers for all request
	resty.SetHeader("Accept", "application/json")
	resty.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"User-Agent":    "goops",
		"Private-Token": viper.GetString("ci_gitlab_token"),
	})
}

func validate() {
	utils.ViperValidateEnv("GOOPSC_GERRIT_URL")
}

func get(endpoint string, response interface{}) {
	validate()
	res, err := resty.R().Get(endpoint)
	if err != nil {
		logrus.Fatalln(err)
	}

	if res.StatusCode() >= 400 {
		logrus.Fatalf("GET: %s\nStatus code: %d\nResponse: %s\n", endpoint, res.StatusCode(), string(res.Body()))
	}

	jsonErr := json.Unmarshal(res.Body()[4:], response)

	if jsonErr != nil {
		logrus.Fatalf("GET: %s\nStatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", endpoint, res.StatusCode(), jsonErr, string(res.Body()[4:]))
	}
}

func GetFiles(changeId string) {
	var response map[string]interface{}
	get(fmt.Sprintf("/changes/%s/revisions/current/files/", changeId), &response)
	for k, _ := range response {
		fmt.Println(k)
	}
}
