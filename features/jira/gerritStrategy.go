package jira

import (
	"github.com/sotomskir/goops/gitService"
	"regexp"
)

type GerritStrategy struct{}

func (GerritStrategy) GetIssues() []string {
	msg := gitService.GetCommitMsg()
	regex := regexp.MustCompile("(\\w+-\\d+)")
	return regex.FindAllString(msg, 99)
}
