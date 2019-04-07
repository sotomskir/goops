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

package gitlabApi

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/utils"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"regexp"
	"time"
)

type MergeRequest struct {
	Id          int    `json:"id,omitempty"`
	Iid         int    `json:"iid,omitempty"`
	ProjectId   int    `json:"project_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
}

type Project struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func Initialize() {
	resty.SetHostURL(viper.GetString("ci_api_v4_url"))
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
	utils.ViperValidate("ci_gitlab_token", "token", "CI_GITLAB_TOKEN")
	utils.ViperValidate("ci_api_v4_url", "server", "CI_API_V4_URL")
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

	jsonErr := json.Unmarshal(res.Body(), response)

	if jsonErr != nil {
		logrus.Fatalf("GET: %s\nStatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", endpoint, res.StatusCode(), jsonErr, string(res.Body()))
	}
}

func post(endpoint string, payload interface{}, response interface{}) {
	validate()
	res, err := resty.R().SetBody(payload).Post(endpoint)
	if err != nil {
		logrus.Fatalln(err)
	}
	if res.StatusCode() >= 400 {
		logrus.Fatalf("POST: %s\nStatus code: %d\nRequest: %#v\nResponse: %s\n", endpoint, res.StatusCode(), payload, string(res.Body()))
	}

	jsonErr := json.Unmarshal(res.Body(), response)
	if jsonErr != nil {
		logrus.Fatalf("POST: %s\nStatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", endpoint, res.StatusCode(), jsonErr, string(res.Body()))
	}
}

func put(endpoint string, payload interface{}, response interface{}) {
	validate()
	res, err := resty.R().SetBody(payload).Put(endpoint)
	if err != nil {
		logrus.Fatalln(err)
	}
	if res.StatusCode() >= 400 {
		logrus.Fatalf("PUT: %s\nStatus code: %d\nRequest: %#v\nResponse: %s\n", endpoint, res.StatusCode(), payload, string(res.Body()))
	}
	if res.StatusCode() == 204 {
		return
	}
	jsonErr := json.Unmarshal(res.Body(), response)
	if jsonErr != nil {
		logrus.Fatalf("PUT: %s\nStatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", endpoint, res.StatusCode(), jsonErr, string(res.Body()))
	}
}

func GetMergeRequestIssueKeys(projectId string, mergeRequestIId string) []string {
	// TODO read keys from commit messages
	mergeRequest := GetMergeRequest(projectId, mergeRequestIId)
	titleKeys := ExtractIssueKeys(mergeRequest.Title)
	descriptionKeys := ExtractIssueKeys(mergeRequest.Description)
	return append(titleKeys, descriptionKeys...)
}

func ExtractIssueKeys(s string) []string {
	regex := regexp.MustCompile("\\w+-\\d+")
	keys := make([]string, 0)
	match := regex.FindAllStringSubmatch(s, -1)
	for _, v := range match {
		keys = append(keys, v[0])
	}
	return keys
}

func GetMergeRequest(projectId string, mergeRequestIId string) MergeRequest {
	mergeRequest := MergeRequest{}
	get(fmt.Sprintf("/projects/%s/merge_requests/%s", projectId, mergeRequestIId), &mergeRequest)
	return mergeRequest
}
