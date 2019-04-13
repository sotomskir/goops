package jira

import (
	"github.com/golang/mock/gomock"
	"github.com/sotomskir/goops/gitService"
	"github.com/sotomskir/goops/mockExecService"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

func TestGetIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tables := []struct {
		msg    string
		issues string
	}{
		{"TEST-1 TEST-2 TEST-321 sad", "TEST-1 TEST-2 TEST-321"},
	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git --no-pager log -1 --pretty=%B").Return(table.msg, nil).AnyTimes()
		gitService.Initialize(mockIService)
		viper.Set("GOOPSC_JIRA", "true")
		j := New()
		actual := j.GetIssues()
		if strings.Join(actual, " ") != table.issues {
			t.Errorf("got: '%s', want: '%s'", strings.Join(actual, " "), table.issues)
		}
	}
}
