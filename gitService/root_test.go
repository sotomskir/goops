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
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/sotomskir/goops/mockExecService"
	"testing"
)

func TestGetHeadTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIService := mock_execService.NewMockIService(ctrl)
	Initialize(mockIService)

	tables := []struct {
		tag      string
		error    error
		expected string
	}{
		{"1.0.0", nil, "1.0.0"},
		{"", nil, ""},
		{"1", errors.New("some error"), ""},
		{"nightly", nil, ""},
	}

	for _, table := range tables {
		mockIService.EXPECT().Exec("git --no-pager tag --contains").Return(table.tag, table.error)
		actual := GetHeadTag()
		if actual != table.expected {
			t.Errorf("Tag is invalid, got: %s, want: %s.", actual, table.expected)
		}
	}
}

func TestGetPreviousTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tables := []struct {
		tag      string
		error    error
		expected string
	}{
		{"1.0.0", nil, "1.0.0"},
		{"", nil, ""},
		{"1", errors.New("some error"), ""},
	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git describe --abbrev=0 --tags --exclude nightly").Return(table.tag, table.error)
		Initialize(mockIService)
		actual := GetPreviousTag()
		if actual != table.expected {
			t.Errorf("Tag is invalid, got: %s, want: %s.", actual, table.expected)
		}
	}
}

func TestStableBranchExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tables := []struct {
		major        int
		minor        int
		stableBranch string
		gitResponse  string
		expected     bool
	}{
		{1, 1, "1.1-stable", "remotes/origin/1.1-stable", true},
		{2, 112, "2.112-stable", "remotes/origin/2.112-stable", true},
		{3, 333, "3.333-stable", "", false},
		{4, 4, "4.4-stable", "", false},
	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec(fmt.Sprintf("git --no-pager branch --remotes --list '*%s'", table.stableBranch)).Return(table.gitResponse, nil).AnyTimes()
		Initialize(mockIService)
		actual := StableBranchExists(table.major, table.minor)
		if actual != table.expected {
			t.Errorf("Version: %d.%d, got: %t, want: %t.", table.major, table.minor, actual, table.expected)
		}
	}
}

func TestGetPreviouslyMergedVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tables := []struct {
		msg        string
		version        string
		error error
	}{
		{"Merge branch release-1.2.3 into master", "1.2.3",nil},
		{"Merge branch release-11.222.3333 into master", "11.222.3333",nil},
		{"0.0.0", "0.0.0",nil},
	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git --no-pager log -n 1 --merges").Return(table.msg, table.error).AnyTimes()
		Initialize(mockIService)
		actual, _ := GetPreviouslyMergedVersion()
		if actual != table.version {
			t.Errorf("TestGetPreviouslyMergedVersion: got: '%s', want: '%s'.", actual, table.version)
		}
	}
}
