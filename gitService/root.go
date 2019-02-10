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

package gitService

import (
	"fmt"
	"github.com/sotomskir/gitlab-cli/execService"
	"github.com/sotomskir/gitlab-cli/logger"
	"os"
	"strconv"
	"strings"
)

var service execService.IService

func Initialize(execServ execService.IService) {
	service = execServ
}

func GetHeadTag() string {
	out, err := service.Exec("git --no-pager tag --contains")
	if err != nil {
		return ""
	}
	return out
}

func GetPreviousTag() string {
	out, err := service.Exec("git describe --abbrev=0 --tags")
	if err != nil {
		return ""
	}
	return out
}

func BumpVersion(version string) string {
	splitted := strings.Split(version, ".")
	minor, err := strconv.Atoi(splitted[1])
	if err != nil {
		logger.ErrorF("Error: BumpVersion probably version: '%s' has invalid semver format.: %#v\n", version, err)
		os.Exit(1)
	}
	splitted[1] = fmt.Sprintf("%d", minor+1)
	patch := strings.Split(splitted[2], "-")
	identifier := strings.Join(patch[1:], "-")
	splitted[2] = "0"
	semver := strings.Join(splitted, ".")
	if identifier != "" {
		semver = fmt.Sprintf("%s-%s", semver, identifier)
	}
	return semver
}

func GetSemanticVersion() string {
	headTag := GetHeadTag()
	if headTag != "" {
		return headTag
	}

	previousTag := GetPreviousTag()
	if previousTag != "" {
		return fmt.Sprintf("%s-SNAPSHOT", BumpVersion(previousTag))
	}

	return "0.1.0-SNAPSHOT"
}
