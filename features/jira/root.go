package jira

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/utils"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/viper"
	"strings"
)

var strategy Strategy

type Strategy interface {
	GetIssues() []string
}

type Jira struct {
	Strategy Strategy
}

func (o *Jira) GetVersion() []string {
	return o.Strategy.GetIssues()
}

func SetJiraVersion(version string, issues []string, summary string, description string, issueType string) {
	if isDisabled() || viper.GetString("GOOPSC_JIRA_VERSION_ASSIGN") == "false" {
		return
	}
	jiraInitApi()
	for _, issue := range issues {
		logrus.Infof("Set version: %s for issue: %s\n", version, issue)
		err := jiraApi.AssignVersion(
			issue,
			version,
			viper.GetString("GOOPSC_JIRA_VERSION_CREATE") != "false",
			viper.GetString("GOOPSC_JIRA_CREATE_DEPLOYMENT_ISSUE") != "false",
			summary,
			description,
			issueType)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

func jiraInit() {
	if isDisabled() {
		return
	}
	viper.SetDefault("GOOPSC_JIRA_STRATEGY", "gerrit")
	switch viper.GetString("GOOPSC_JIRA_STRATEGY") {
	case "gerrit":
		strategy = GerritStrategy{}
		break
	case "gitlab":
		strategy = GitlabStrategy{}
		break
	default:
		panic(fmt.Sprintf("unsupported strategy: %s\n", viper.GetString("GOOPSC_JIRA_STRATEGY")))
	}
}

func jiraInitApi() {
	utils.ViperValidateEnv("GOOPSC_JIRA_SERVER_URL", "GOOPSC_JIRA_USER", "GOOPSC_JIRA_PASSWORD")
	jiraApi.Initialize(viper.GetString("GOOPSC_JIRA_SERVER_URL"), viper.GetString("GOOPSC_JIRA_USER"), viper.GetString("GOOPSC_JIRA_PASSWORD"))
}

func isDisabled() bool {
	viper.SetDefault("GOOPSC_JIRA", "false")
	return viper.GetString("GOOPSC_JIRA") == "false"
}

func GetIssues() []string {
	if isDisabled() {
		return nil
	}
	jiraInit()
	issues := strategy.GetIssues()
	issueKeysJoined := strings.Join(issues, " ")
	utils.SaveExportString("GOOPS_ISSUES", issueKeysJoined)
	return issues
}

func JiraTransition(issues string, state string) {
	if isDisabled() || viper.GetString("GOOPSC_JIRA_ISSUE_TRANSITION") == "false" {
		return
	}
	jiraInitApi()
	for _, issue := range strings.Split(issues, " ") {
		logrus.Infof("Transition issue: %s to state: %s\n", issue, state)
		jiraApi.TransitionIssue("", issue, state, "")
	}
}
