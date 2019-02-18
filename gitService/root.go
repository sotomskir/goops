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
	major, minor, _, identifier := SplitSemver(version)
	return fmt.Sprintf("%d.%d.%d%s", major, minor + 1, 0, identifier)
}

func BumpPatchVersion(version string) string {
	major, minor, patch, identifier := SplitSemver(version)
	return fmt.Sprintf("%d.%d.%d%s", major, minor, patch + 1, identifier)
}

func SplitSemver(version string) (int, int, int, string) {
	splited := strings.Split(version, ".")
	if len(splited) < 3 {
		logrus.Fatalf("Invalid semver format: %s\n", version)
	}
	major, err := strconv.Atoi(splited[0])
	if err != nil {
		logrus.Fatalf("Error converting major version: %s to int\n", splited[0])
	}
	minor, err := strconv.Atoi(splited[1])
	if err != nil {
		logrus.Fatalf("Error converting patch version: %s to int\n", splited[1])
	}
	splitedPatch := strings.Split(splited[2], "-")
	patch, err := strconv.Atoi(splitedPatch[0])
	if err != nil {
		logrus.Fatalf("Error converting patch version: %s to int\n", splitedPatch[0])
	}
	identifier := strings.Join(splitedPatch[1:], "-")
	if identifier != "" {
		identifier = "-" + identifier
	}
	return major, minor, patch, identifier
}

func GetSemanticVersion() (string, string) {
	headTag := GetHeadTag()
	if headTag != "" {
		return headTag, headTag
	}
	previousTag := GetPreviousTag()
	if previousTag == "" {
		return "0.1.0-SNAPSHOT", "0.1.0"
	}
	var version string
	if IsStableBranch(GetCurrentBranchName()) {
		version = GetVersionForStableBranch(previousTag)
	} else {
		version = BumpMinorVersion(previousTag)
		println(version)
		if StableBranchExists(version) {
			version = BumpMinorVersion(version)
		}
	}
	return fmt.Sprintf("%s-SNAPSHOT", version), version
}

func StableBranchExists(version string) bool {
	major, minor, _, _ := SplitSemver(version)
	res, err := service.Exec(fmt.Sprintf("git --no-pager branch --remotes --list '*%d.%d-stable'", major, minor))
	if err != nil {
		logrus.Fatalln(res, err)
	}
	return res != ""
}

func GetVersionForStableBranch(previousTag string) string {
	if VersionMatchBranchName(previousTag, GetCurrentBranchName()) {
		return BumpPatchVersion(previousTag)
	}
	return GetVersionFromBranchName(GetCurrentBranchName())
}

func GetVersionFromBranchName(branch string) string {
	if IsStableBranch(branch) {
		regex := regexp.MustCompile("^(.*)-stable")
		match := regex.FindStringSubmatch(branch)
		if len(match) < 2 {
			logrus.Fatalf("Version not found in branch name: %s\n", branch)
		}
		return fmt.Sprintf("%s.0", match[1])
	}
	regex := regexp.MustCompile("(\\d+\\.\\d+\\.\\d+)")
	match := regex.FindStringSubmatch(branch)
	if len(match) < 2 {
		logrus.Fatalf("Version not found in branch name: %s\n", branch)
	}
	return match[1]
}

func VersionMatchBranchName(version string, branch string) bool {
	regex := regexp.MustCompile("(\\d+\\.\\d+)")
	match := regex.FindStringSubmatch(branch)
	match2 := regex.FindStringSubmatch(version)
	if len(match) < 2 || len(match2) < 2 {
		logrus.Fatalf("Version not found in branch name: %s\n", branch)
	}
	return match[1] == match2[1]
}

func IsStableBranch(branch string) bool {
	match, _ := regexp.MatchString("^.*-stable$", branch)
	match2, _ := regexp.MatchString("^.*-stable$", viper.GetString("CI_COMMIT_REF_NAME"))
	return match || match2
}

func GetCurrentBranchName() string {
	branch, err := service.Exec("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		logrus.Fatalln(branch, err)
	}
	return branch
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
