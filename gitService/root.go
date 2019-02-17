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
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/gitlab-cli/execService"
	"github.com/spf13/viper"
	"regexp"
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

func BumpMinorVersion(version string) string {
	splitted := strings.Split(version, ".")
	minor, err := strconv.Atoi(splitted[1])
	if err != nil {
		logrus.Fatalf("Error: BumpMinorVersion probably version: '%s' has invalid semver format.: %#v\n", version, err)
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

func BumpPatchVersion(version string) string {
	splited := strings.Split(version, ".")
	major := splited[0]
	minor := splited[1]
	patch, err := strconv.Atoi(strings.Split(splited[2], "-")[0])
	if err != nil {
		logrus.Fatalf("Error: BumpPatchVersion probably version: '%s' has invalid semver format.: %#v\n", version, err)
	}
	identifier := strings.Join(strings.Split(splited[2], "-")[1:], "-")
	splited[2] = fmt.Sprintf("%d", patch+1)
	semver := fmt.Sprintf("%s.%s.%d", major, minor, patch + 1)
	if identifier != "" {
		semver = fmt.Sprintf("%s-%s", semver, identifier)
	}
	return semver
}

func GetSemanticVersion() (string, string) {
	headTag := GetHeadTag()
	if headTag != "" {
		return headTag, headTag
	}

	previousTag := GetPreviousTag()
	if previousTag != "" {
		version := previousTag
		if IsStableBranch() {
			version = BumpPatchVersion(version)
		} else {
			version = BumpMinorVersion(version)
		}
		return fmt.Sprintf("%s-SNAPSHOT", version), version
	}

	return "0.1.0-SNAPSHOT", "0.1.0"
}

func IsStableBranch() bool {
	branch, err := service.Exec("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		logrus.Fatalln(branch, err)
	}
	match, _ := regexp.MatchString("^.*-stable$", branch)
	match2, _ := regexp.MatchString("^.*-stable$", viper.GetString("CI_COMMIT_REF_NAME"))
	return match || match2
}

func GetPreviousMergeRequestIid() string {
	previousMerge, err := service.Exec("git --no-pager log -1 --merges")
	if err != nil {
		logrus.Fatalln(err)
	}
	if previousMerge == "" {
		logrus.Fatalln("Merge request not found")
	}
	return ExtractMergeRequestIid(previousMerge)
}

func ExtractMergeRequestIid(s string) string {
	regex := regexp.MustCompile("!(\\d+)")
	match := regex.FindStringSubmatch(s)
	if len(match) < 2 {
		logrus.Fatalln("Merge request not found")
	}
	return match[1]
}
