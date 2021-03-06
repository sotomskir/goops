package semver

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/sotomskir/goops/gitService"
	"github.com/sotomskir/goops/mockExecService"
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
		actual := bumpMinorVersion(table.version)
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
		actual := bumpPatchVersion(table.version)
		if actual != table.expected {
			t.Errorf("Version is invalid for input: '%s', got: '%s', want: '%s'", table.version, actual, table.expected)
		}
	}
}

func TestIsStableBranch(t *testing.T) {
	tables := []struct {
		branch   string
		expected bool
	}{
		{"master", false},
		{"1.1-stable", true},
		{"112.346-stable", true},
		{"64-123-stable", true},
		{"hotfix-4.2.311", false},
	}

	for _, table := range tables {
		actual := isStableBranch(table.branch)
		if actual != table.expected {
			t.Errorf("branch: %s, got: %t, want: %t.", table.branch, actual, table.expected)
		}
	}

	viper.Set("CI_COMMIT_REF_NAME", "6.7-stable")
	actual := isStableBranch("da7f8adfy7asfd7")
	if actual != true {
		t.Errorf("branch: %s, got: %t, want: %t.", "da7f8adfy7asfd7", actual, true)
	}
	viper.Set("CI_COMMIT_REF_NAME", "")
}

func TestGetVersionFromBranchName(t *testing.T) {
	tables := []struct {
		branch   string
		expected string
	}{
		{"1.0-stable", "1.0.0"},
		{"2.11-stable", "2.11.0"},
		{"hotfix-3.2.311", "3.2.311"},
		{"release-4.51.1", "4.51.1"},
	}

	for _, table := range tables {
		actual := getVersionFromBranchName(table.branch)
		if actual != table.expected {
			t.Errorf("branch: %s, got: %s, want: %s.", table.branch, actual, table.expected)
		}
	}
}

func TestVersionMatchBranchName(t *testing.T) {
	tables := []struct {
		version  string
		branch   string
		expected bool
	}{
		{"1.9.0", "2.0-stable", false},
		{"2.99.0", "2.99-stable", true},
		{"3.0.23", "3.0-stable", true},
		{"4.44.4", "4.45-stable", false},
	}

	for _, table := range tables {
		actual := versionMatchBranchName(table.version, table.branch)
		if actual != table.expected {
			t.Errorf("version: %s, branch: %s, got: %t, want: %t.", table.version, table.branch, actual, table.expected)
		}
	}
}

func TestSplitSemver(tst *testing.T) {
	tables := []struct {
		version    string
		major      int
		minor      int
		patch      int
		identifier string
	}{
		{"1.9.0", 1, 9, 0, ""},
		{"2.10.99", 2, 10, 99, ""},
		{"3.909.1220-SNAPSHOT", 3, 909, 1220, "-SNAPSHOT"},
	}

	for _, t := range tables {
		major, minor, patch, identifier := splitSemver(t.version)
		if major != t.major || minor != t.minor || patch != t.patch || identifier != t.identifier {
			tst.Errorf("version: %s, got: %d.%d.%d%s, want: %d.%d.%d%s", t.version, major, minor, patch, identifier, t.major, t.minor, t.patch, t.identifier)
		}
	}
}

func TestGetSemanticVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tables := []struct {
		branch             string
		tag                string
		previousTag        string
		error              error
		previousError      error
		stableBranch       string
		stableBranchReturn string
		expected           string
	}{
		// If current HEAD is tagged then tag will be used as version.
		{"1.1-stable", "1.1.2", "1.1.1", nil, nil, "", "", "1.1.2"},

		// Else command will lookup for previous tag bump it's minor version, reset patch version and append '-SNAPSHOT'
		{"master", "", "1.12.3", nil, nil, "1.13-stable", "", "1.13.0-SNAPSHOT"},
		{"master", "", "4.44.444", nil, nil, "4.45-stable", "remotes/origin/4.45-stable", "4.46.0-SNAPSHOT"},

		// If branch is *-stable will lookup for previous tag and bump it's patch version
		{"2.2-stable", "", "2.2.12", nil, nil, "", "", "2.2.13-SNAPSHOT"},

		// When there are no tags found version will be '0.1.0-SNAPSHOT'`,
		{"master", "", "", nil, errors.New("nothing to describe"), "0.1-stable", "", "0.1.0-SNAPSHOT"},

		// If branch is *-stable and previous tag is not from current stable version, will take version from branch name and set patch to 0
		{"3.10-stable", "", "3.9.12", nil, nil, "", "", "3.10.0-SNAPSHOT"},

		{"master", "", "", nil, nil, "0.1-stable", "", "0.1.0-SNAPSHOT"},
		{"master", "", "", nil, nil, "0.1-stable", "remotes/origin/0.1-stable", "0.2.0-SNAPSHOT"},
		{"master", "", "0.1.0", nil, nil, "0.2-stable", "", "0.2.0-SNAPSHOT"},
		{"master", "", "0.1.1", nil, nil, "0.2-stable", "", "0.2.0-SNAPSHOT"},
		{"master", "", "0.1.0", nil, nil, "0.2-stable", "", "0.2.0-SNAPSHOT"},
		{"0.1-stable", "", "", nil, nil, "0.1-stable", "remotes/origin/0.1-stable", "0.1.0-SNAPSHOT"},
		{"0.1-stable", "0.1.0", "", nil, nil, "0.1-stable", "remotes/origin/0.1-stable", "0.1.0"},
		{"0.1-stable", "", "0.1.0", nil, nil, "0.1-stable", "remotes/origin/0.1-stable", "0.1.1-SNAPSHOT"},
		{"0.1-stable", "0.1.1", "0.1.0", nil, nil, "0.1-stable", "remotes/origin/0.1-stable", "0.1.1"},
		{"0.2-stable", "", "0.1.0", nil, nil, "0.2-stable", "remotes/origin/0.2-stable", "0.2.0-SNAPSHOT"},
	}

	viper.Set(GoopscSemverSaveExport, "false")
	viper.Set(GoopscSemver, "true")
	viper.Set(GoopscSemverStrategy, GitlabFlowStrategy)
	s := New()
	for _, table := range tables {
		mockIService := mock_execService.NewMockIService(ctrl)
		mockIService.EXPECT().Exec("git --no-pager tag --contains").Return(table.tag, table.error).AnyTimes()
		mockIService.EXPECT().Exec("git describe --abbrev=0 --tags --exclude nightly").Return(table.previousTag, table.previousError).AnyTimes()
		mockIService.EXPECT().Exec("git rev-parse --abbrev-ref HEAD").Return(table.branch, nil).AnyTimes()
		mockIService.EXPECT().Exec(fmt.Sprintf("git --no-pager branch --remotes --list '*%s'", table.stableBranch)).Return(table.stableBranchReturn, nil).AnyTimes()
		gitService.Initialize(mockIService)
		actual := s.GetVersion()
		if actual != table.expected {
			t.Errorf("Version is invalid, got: '%s', want: '%s'\n%v.", actual, table.expected, table)
		}
	}
}
