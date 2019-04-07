package jira

import (
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/gitService"
	"github.com/sotomskir/goops/gitlabApi"
	"github.com/spf13/viper"
)

type GitlabStrategy struct{}

func (GitlabStrategy) GetIssues() []string {
	mergeRequestIid := getMergeRequestIid()
	mergeRequestIid = "3"
	projectId := viper.GetString("CI_PROJECT_ID")
	if projectId == "" {
		logrus.Fatalln("CI_PROJECT_ID is not set")
	}
	issueKeys := gitlabApi.GetMergeRequestIssueKeys(projectId, mergeRequestIid)
	return issueKeys
}

func getMergeRequestIid() string {
	iid := viper.GetString("CI_MERGE_REQUEST_IID")
	if iid == "" {
		iid = gitService.GetPreviousMergeRequestIid()
	}
	return iid
}
