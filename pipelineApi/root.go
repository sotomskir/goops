package pipelineApi

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/gitlab-cli/execService"
	"github.com/sotomskir/gitlab-cli/gitService"
	"github.com/sotomskir/gitlab-cli/gitlabApi"
	"github.com/sotomskir/gitlab-cli/utils"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"
)

var e execService.IService

func Initialize(execService execService.IService) {
	e = execService
}

func Common() {
	version := semanticVersion()
	issues := issues()
	jiraVersion(version, issues)
}

func jiraVersion(version string, issues []string) {
	jiraInit()
	for _, issue := range issues {
		logrus.Infof("Set version: %s for issue: %s\n", version, issue)
		jiraApi.SetFixVersion(issue, version)
	}
}

func jiraInit() {
	utils.ViperValidateEnv("JIRA_SERVER_URL", "JIRA_USER", "JIRA_PASSWORD")
	jiraApi.Initialize(viper.GetString("JIRA_SERVER_URL"), viper.GetString("JIRA_USER"), viper.GetString("JIRA_PASSWORD"))
}

func issues() []string {
	mergeRequestIid := getMergeRequestIid()
	mergeRequestIid = "3"
	projectId := viper.GetString("CI_PROJECT_ID")
	if projectId == "" {
		logrus.Fatalln("CI_PROJECT_ID is not set")
	}
	issueKeys := gitlabApi.GetMergeRequestIssueKeys(projectId, mergeRequestIid)
	issueKeysJoined := strings.Join(issueKeys, " ")
	saveExport("CI_ISSUES", issueKeysJoined)
	return issueKeys
}

func semanticVersion() string {
	semver, releaseSemver := gitService.GetSemanticVersion()
	saveExport("CI_SEMVER", semver)
	saveExport("CI_SEMVER_RELEASE", releaseSemver)
	return releaseSemver
}

func saveExport(variableName string, value string) {
	f, err := os.OpenFile("gitlab.env", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Fatalln(err)
	}

	defer f.Close()

	s := fmt.Sprintf("export %s=\"%s\"\n", variableName, value)
	logrus.Infoln(s)
	if _, err = f.WriteString(s); err != nil {
		logrus.Fatalln(err)
	}
}

func getMergeRequestIid() string {
	iid := viper.GetString("CI_MERGE_REQUEST_IID")
	if iid == "" {
		iid = gitService.GetPreviousMergeRequestIid()
	}
	return iid
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

func JiraTransition(issues string, state string) {
	jiraInit()
	for _, issue := range strings.Split(issues, " ") {
		logrus.Infof("Transition issue: %s to state: %s\n", issue, state)
		jiraApi.TransitionIssue("", issue, state)
	}
}
