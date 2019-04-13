package jira

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/utils"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/viper"
	"strings"
)

const (
	// Configuration variables
	GoopscJira                      = "GOOPSC_JIRA"
	GoopscJiraProjectKey            = "GOOPSC_JIRA_PROJECT_KEY"
	GoopscJiraServerUrl             = "GOOPSC_JIRA_SERVER_URL"
	GoopscJiraUser                  = "GOOPSC_JIRA_USER"
	GoopscJiraPassword              = "GOOPSC_JIRA_PASSWORD"
	GoopscJiraVersionAssign         = "GOOPSC_JIRA_VERSION_ASSIGN"
	GoopscJiraVersionCreate         = "GOOPSC_JIRA_VERSION_CREATE"
	GoopscJiraCreateDeploymentIssue = "GOOPSC_JIRA_CREATE_DEPLOYMENT_ISSUE"
	GoopscJiraIssueTransition       = "GOOPSC_JIRA_ISSUE_TRANSITION"
	GoopscJiraWorkflow              = "GOOPSC_JIRA_WORKFLOW"
	GoopscJiraWorkflowContent       = "GOOPSC_JIRA_WORKFLOW_CONTENT"
	GoopscJiraStrategy              = "GOOPSC_JIRA_STRATEGY"

	// Output variables
	GoopsJiraIssues = "GOOPS_JIRA_ISSUES"

	// Configuration options
	GerritStrategy = "gerrit"
	GitlabStrategy = "gitlab"
)

func setDefaults() {
	viper.SetDefault(GoopscJira, "false")
	viper.SetDefault(GoopscJiraStrategy, GerritStrategy)
}

type strategy interface {
	getIssues() []string
}

type Jira struct {
	strategy strategy
}

func New() Jira {
	setDefaults()
	var strategy strategy
	switch viper.GetString(GoopscJiraStrategy) {
	case GerritStrategy:
		strategy = gerritStrategy{}
		break
	case GitlabStrategy:
		strategy = gitlabStrategy{}
		break
	default:
		panic(fmt.Sprintf("unsupported strategy: %s\n", viper.GetString(GoopscJiraStrategy)))
	}
	return Jira{strategy: strategy}
}

func (o *Jira) GetIssues() []string {
	if utils.IsDisabled(GoopscJira) {
		return nil
	}
	issues := o.strategy.getIssues()
	issueKeysJoined := strings.Join(issues, " ")
	utils.SaveExportString(GoopsJiraIssues, issueKeysJoined)
	return issues
}

func (o *Jira) JiraTransition(issues string, state string) {
	if utils.IsDisabled(GoopscJira) || utils.IsDisabled(GoopscJiraIssueTransition) {
		return
	}
	jiraInitApi()
	for _, issue := range strings.Split(issues, " ") {
		logrus.Infof("Transition issue: %s to state: %s\n", issue, state)
		jiraApi.TransitionIssue("", issue, state, "")
	}
}

func (o *Jira) SetJiraVersion(version string, issues []string, summary string, description string, issueType string) {
	if utils.IsDisabled(GoopscJira) || utils.IsDisabled(GoopscJiraVersionAssign) {
		return
	}
	jiraInitApi()
	for _, issue := range issues {
		logrus.Infof("Set version: %s for issue: %s\n", version, issue)
		err := jiraApi.AssignVersion(
			issue,
			version,
			utils.IsEnabled(GoopscJiraVersionCreate),
			utils.IsEnabled(GoopscJiraCreateDeploymentIssue),
			summary,
			description,
			issueType)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

func jiraInitApi() {
	utils.ViperValidateEnv(GoopscJiraServerUrl, GoopscJiraUser, GoopscJiraPassword)
	jiraApi.Initialize(viper.GetString(GoopscJiraServerUrl), viper.GetString(GoopscJiraUser), viper.GetString(GoopscJiraPassword))
}
