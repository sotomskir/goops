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
	"github.com/golang/mock/gomock"
	"github.com/sotomskir/gitlab-cli/mockExecService"
	"testing"
)

func TestBumpVersion(t *testing.T) {
	tables := []struct {
		version  string
		expected string
	}{
		{"1.0.1", "1.1.0"},
		{"1.111.0", "1.112.0"},
		{"0.0.123", "0.1.0"},
		{"1123.1212.0-SNAPSHOT", "1123.1213.0-SNAPSHOT"},
		{"1123.1212.12-SNAPSHOT", "1123.1213.0-SNAPSHOT"},
	}

	for _, table := range tables {
		actual := BumpVersion(table.version)
		if actual != table.expected {
			t.Errorf("Version is invalid for input: '%s', got: '%s', want: '%s'", table.version, actual, table.expected)
		}
	}
}

func TestGetHeadTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIService := mock_execService.NewMockIService(ctrl)
	Initialize(mockIService)

	tables := []struct {
		tag  string
		error error
		expected string
	}{
		{"1.0.0", nil, "1.0.0"},
		{"", nil, ""},
		{"1", errors.New("some error"), ""},
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
		tag  string
		error error
		expected string
	}{
		{"1.0.0", nil, "1.0.0"},
		{"", nil, ""},
		{"1", errors.New("some error"), ""},
	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git describe --abbrev=0 --tags").Return(table.tag, table.error)
		Initialize(mockIService)
		actual := GetPreviousTag()
		if actual != table.expected {
			t.Errorf("Tag is invalid, got: %s, want: %s.", actual, table.expected)
		}
	}
}

func TestGetSemanticVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tables := []struct {
		tag  string
		previousTag  string
		error error
		previousError error
		expected string
	}{
		// If current HEAD is tagged then tag will be used as version.
		{"1.1.2", "1.1.1", nil, nil, "1.1.2"},

	   	// Else command will lookup for previous tag bump it's minor version, reset patch version and append '-SNAPSHOT'
		{"", "1.12.3", nil, nil, "1.13.0-SNAPSHOT"},

		// When there are no tags found version will be '0.1.0-SNAPSHOT'`,
		{"", "", nil, errors.New("nothing to describe"), "0.1.0-SNAPSHOT"},
	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git --no-pager tag --contains").Return(table.tag, table.error).AnyTimes()
		mockIService.EXPECT().Exec("git describe --abbrev=0 --tags").Return(table.previousTag, table.previousError).AnyTimes()
		Initialize(mockIService)
		actual := GetSemanticVersion()
		if actual != table.expected {
			t.Errorf("Version is invalid, got: %s, want: %s.", actual, table.expected)
		}
	}
}
