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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/execService"
	"github.com/spf13/viper"
	"regexp"
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
	return strings.Trim(out, " \n\t")
}

func GetPreviousTag() string {
	out, err := service.Exec("git describe --abbrev=0 --tags")
	if err != nil {
		return ""
	}
	return strings.Trim(out, " \n\t")
}

func StableBranchExists(major int, minor int) bool {
	res, err := service.Exec(fmt.Sprintf("git --no-pager branch --remotes --list '*%d.%d-stable'", major, minor))
	if err != nil {
		logrus.Fatalln(res, err)
	}
	return res != ""
}

func BranchExists(version string) bool {
	res, err := service.Exec(fmt.Sprintf("git --no-pager branch --remotes --list '*%s*'", version))
	if err != nil {
		logrus.Fatalln(res, err)
	}
	return res != ""
}

func GetCurrentBranchName() string {
	branch, err := service.Exec("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		logrus.Fatalln(branch, err)
	}
	return branch
}

func GetCommitMsg() string {
	msg, err := service.Exec("git --no-pager log -1 --pretty=%B")
	if err != nil {
		logrus.Fatalln(msg, err)
	}
	return msg
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

func setupGit() {
	viper.SetDefault("GOOPSC_GIT_USER_EMAIL", "travis@travis-ci.org")
	viper.SetDefault("GOOPSC_GIT_USER_NAME", "Travis CI")
	service.Exec(fmt.Sprintf("git config --global user.email '%s'", viper.GetString("GOOPSC_GIT_USER_EMAIL")))
	service.Exec(fmt.Sprintf("git config --global user.name '%s'", viper.GetString("GOOPSC_GIT_USER_NAME")))
}

func TagNightly() {
	if viper.GetString("TRAVIS_PULL_REQUEST") == "false" && viper.GetString("TRAVIS_BRANCH") == "master" && viper.GetString("TRAVIS_TAG") == "" {
		setupGit()
		service.LogExec("git tag -f nightly")
		output, err := service.Exec("git --no-pager remote")
		if err != nil {
			panic(err)
		}
		if strings.Contains(output, "goops-remote") {
			service.LogExec("git remote remove goops-remote")
		}
		_, err = service.Exec(
			fmt.Sprintf(
				"git remote add goops-remote https://%s:%s@github.com/%s.git",
				viper.GetString("GOOPSC_GITHUB_USER"),
				viper.GetString("GOOPSC_GITHUB_TOKEN"),
				viper.GetString("TRAVIS_REPO_SLUG")))
		if err != nil {
			panic("error adding remote goops-remote")
		}
		_, err = service.Exec("git push -f --tags goops-remote")
		if err != nil {
			panic("error pushing tag")
		}
		service.LogExec("git remote remove goops-remote")
	}
}

func GetPreviouslyMergedVersion() (string, error) {
	msg, err := service.Exec("git --no-pager log -n 1 --merges")
	if err != nil {
		return "", err
	}
	if msg == "" {
		msg = "0.0.0"
	}
	regex := regexp.MustCompile("(\\d+\\.\\d+\\.\\d+)")
	match := regex.FindStringSubmatch(msg)
	if len(match) < 1 {
		return "", errors.New("Version not found")
	}
	return match[0], nil
}
