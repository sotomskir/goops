package jira

import (
	"github.com/sotomskir/goops/gitService"
	"regexp"
)

type gerritStrategy struct{}

func (gerritStrategy) getIssues() []string {
	msg := gitService.GetCommitMsg()
	regex := regexp.MustCompile("(\\w+-\\d+)")
	return regex.FindAllString(msg, 99)
}
