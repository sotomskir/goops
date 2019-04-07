package docker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/execService"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

var e execService.IService

func Initialize(execService execService.IService) {
	e = execService
}

func DockerBuild(tag string, dockerfile string, path string) {
	if dockerfile != "" {
		dockerfile = fmt.Sprintf("-f %s", dockerfile)
	}
	if tag != "" {
		tag = fmt.Sprintf("-t %s", tag)
	}
	e.LogExec(fmt.Sprintf("docker build %s %s %s", dockerfile, tag, path))
}

func DockerPush(tag string) {
	refName := viper.GetString("CI_COMMIT_REF_NAME")
	isMaster := refName == "master"
	isTag := viper.GetString("CI_COMMIT_TAG") != ""
	isStable, _ := regexp.MatchString("^.*-stable$", refName)
	if !isMaster && !isStable && !isTag {
		logrus.Infoln("Docker publish skipped")
		logrus.Debugf("refName=\"%s\", tag=\"%s\"\n", refName, viper.GetString("CI_COMMIT_TAG"))
		return
	}
	if isTag {
		image := strings.Split(tag, ":")[0]
		e.LogExec(fmt.Sprintf("docker tag %s %s:stable", tag, image))
		e.LogExec(fmt.Sprintf("docker push %s:stable", image))
	}
	if isMaster && !isTag {
		image := strings.Split(tag, ":")[0]
		e.LogExec(fmt.Sprintf("docker tag %s %s:latest", tag, image))
		e.LogExec(fmt.Sprintf("docker push %s:latest", image))
	}
	e.LogExec(fmt.Sprintf("docker push %s", tag))
}
