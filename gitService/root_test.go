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
	"github.com/spf13/viper"
	"testing"
)

func TestBumpMinorVersion(t *testing.T) {
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
		actual := BumpMinorVersion(table.version)
		if actual != table.expected {
			t.Errorf("Version is invalid for input: '%s', got: '%s', want: '%s'", table.version, actual, table.expected)
		}
	}
}

func TestBumpPatchVersion(t *testing.T) {
	tables := []struct {
		version  string
		expected string
	}{
		{"1.0.1", "1.0.2"},
		{"1.111.0", "1.111.1"},
		{"0.0.123", "0.0.124"},
		{"1123.1212.0-SNAPSHOT", "1123.1212.1-SNAPSHOT"},
		{"1123.1212.12-SNAPSHOT", "1123.1212.13-SNAPSHOT"},
	}

	for _, table := range tables {
		actual := BumpPatchVersion(table.version)
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
		branch  string
		tag  string
		previousTag  string
		error error
		previousError error
		expected string
		expectedRelease string
	}{
		// If current HEAD is tagged then tag will be used as version.
		{"1.1-stable", "1.1.2", "1.1.1", nil, nil, "1.1.2", "1.1.2"},

	   	// Else command will lookup for previous tag bump it's minor version, reset patch version and append '-SNAPSHOT'
		{"master", "", "1.12.3", nil, nil, "1.13.0-SNAPSHOT", "1.13.0"},

		// If branch is *-stable will lookup for previous tag and bump it's patch version
		{"2.2-stable", "", "2.2.12", nil, nil, "2.2.13-SNAPSHOT", "2.2.13"},

		// When there are no tags found version will be '0.1.0-SNAPSHOT'`,
		{"master", "", "", nil, errors.New("nothing to describe"), "0.1.0-SNAPSHOT", "0.1.0"},

	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git --no-pager tag --contains").Return(table.tag, table.error).AnyTimes()
		mockIService.EXPECT().Exec("git describe --abbrev=0 --tags").Return(table.previousTag, table.previousError).AnyTimes()
		mockIService.EXPECT().Exec("git rev-parse --abbrev-ref HEAD").Return(table.branch, nil).AnyTimes()
		Initialize(mockIService)
		actual, actualRelease := GetSemanticVersion()
		if actual != table.expected {
			t.Errorf("Version is invalid, got: %s, want: %s.", actual, table.expected)
		}
		if actualRelease != table.expectedRelease {
			t.Errorf("Release version is invalid, got: %s, want: %s.", actualRelease, table.expectedRelease)
		}
	}
}

func TestIsStableBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tables := []struct {
		branch  string
		expected bool
	}{
		{"master", false},
		{"1.1-stable", true},
		{"112.346-stable", true},
		{"64-123-stable", true},

	}

	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git rev-parse --abbrev-ref HEAD").Return(table.branch, nil).AnyTimes()
		Initialize(mockIService)
		actual := IsStableBranch()
		if actual != table.expected {
			t.Errorf("branch: %s, got: %t, want: %t.", table.branch, actual, table.expected)
		}
	}

	mockIService := mock_execService.NewMockIService(ctrl)
	mockIService.EXPECT().Exec("git rev-parse --abbrev-ref HEAD").Return("da7f8adfy7asfd7", nil).AnyTimes()
	Initialize(mockIService)
	viper.Set("CI_COMMIT_REF_NAME", "6.7-stable")
	actual := IsStableBranch()
	if actual != true {
		t.Errorf("branch: %s, got: %t, want: %t.", "da7f8adfy7asfd7", actual, true)
	}

}
