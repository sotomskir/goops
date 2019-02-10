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
	"github.com/sotomskir/gitlab-cli/logger"
	"github.com/sotomskir/gitlab-cli/utils"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"os"
	"regexp"
	"time"
)

type MergeRequest struct {
	Id int `json:"id,omitempty"`
	Iid int `json:"iid,omitempty"`
	ProjectId int `json:"project_id,omitempty"`
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	State string `json:"state,omitempty"`
}

type Project struct {
	Id  string `json:"id"`
	Key string `json:"key"`
	Name string `json:"name"`
}

func Initialize() {
	utils.ViperValidate("build_token", "token", "CI_BUILD_TOKEN")
	utils.ViperValidate("api_v4_url", "server", "CI_API_V4_URL")
	resty.SetHostURL(viper.GetString("api_v4_url"))
	resty.SetTimeout(1 * time.Minute)

	// Headers for all request
	resty.SetHeader("Accept", "application/json")
	resty.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"User-Agent":    "gitlab-cli",
		"Private-Token": viper.GetString("build_token"),
	})
}

func get(endpoint string, response interface{}) {
	res, err := resty.R().Get(endpoint)
	if err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}

	if res.StatusCode() >= 400 {
		logger.ErrorF("Status code: %d\nResponse: %s\n", res.StatusCode(), string(res.Body()))
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(res.Body(), response)

	if jsonErr != nil {
		logger.ErrorF("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func post(endpoint string, payload interface{}, response interface{}) {
	res, err := resty.R().SetBody(payload).Post(endpoint)
	if err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}
	if res.StatusCode() >= 400 {
		logger.ErrorF("Status code: %d\nRequest: %#v\nResponse: %s\n", res.StatusCode(), payload, string(res.Body()))
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(res.Body(), response)
	if jsonErr != nil {
		logger.ErrorF("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func put(endpoint string, payload interface{}, response interface{}) {
	res, err := resty.R().SetBody(payload).Put(endpoint)
	if err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}
	if res.StatusCode() >= 400 {
		logger.ErrorF("Status code: %d\nRequest: %#v\nResponse: %s\n", res.StatusCode(), payload, string(res.Body()))
		os.Exit(1)
	}
	if res.StatusCode() == 204 {
		return
	}
	jsonErr := json.Unmarshal(res.Body(), response)
	if jsonErr != nil {
		logger.ErrorF("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func GetMergeRequestIssueKeys(projectId string, mergeRequestIId string) []string {
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
