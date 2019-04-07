package semver

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/goops/gitService"
	"github.com/sotomskir/goops/utils"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Strategy interface {
	GetSemanticVersion() string
}

type Semver struct {
	Strategy Strategy
}

func (o *Semver) GetVersion() string {
	return o.Strategy.GetSemanticVersion()
}

func GetSemanticVersion() string {
	if viper.GetString("GOOPSC_SEMVER") == "false" {
		return ""
	}
	var strategy Strategy
	switch viper.GetString("GOOPSC_SEMVER_STRATEGY") {
	case "github-flow":
		strategy = GithubFlow{}
	case "gitlab-flow":
		strategy = GitlabFlow{}
	case "git-flow-branch":
		strategy = GitFlowBranch{}
	case "":
		strategy = GitlabFlow{}
	default:
		logrus.Errorf("Unexpected strategy: %s\n", viper.GetString("GOOPSC_SEMVER_STRATEGY"))
		os.Exit(1)
	}
	s := Semver{Strategy: strategy}
	version := s.GetVersion()
	utils.SaveExportString("GOOPS_SEMVER", version)
	utils.SaveExportString("GOOPS_SEMVER_RELEASE", strings.Replace(version, "-SNAPSHOT", "", 1))
	major, minor, patch, _ := splitSemver(version)
	utils.SaveExportInt("GOOPS_SEMVER_MAJOR", major)
	utils.SaveExportInt("GOOPS_SEMVER_MINOR", minor)
	utils.SaveExportInt("GOOPS_SEMVER_PATCH", patch)
	return version
}

func bumpMinorVersion(version string) string {
	major, minor, _, identifier := splitSemver(version)
	return fmt.Sprintf("%d.%d.%d%s", major, minor+1, 0, identifier)
}

func bumpPatchVersion(version string) string {
	major, minor, patch, identifier := splitSemver(version)
	return fmt.Sprintf("%d.%d.%d%s", major, minor, patch+1, identifier)
}

func splitSemver(version string) (int, int, int, string) {
	splited := strings.Split(version, ".")
	if len(splited) < 3 {
		logrus.Fatalf("Invalid semver format: %s\n", version)
	}
	major, err := strconv.Atoi(splited[0])
	if err != nil {
		logrus.Fatalf("Error converting major version: %s to int\n", splited[0])
	}
	minor, err := strconv.Atoi(splited[1])
	if err != nil {
		logrus.Fatalf("Error converting patch version: %s to int\n", splited[1])
	}
	splitedPatch := strings.Split(splited[2], "-")
	patch, err := strconv.Atoi(splitedPatch[0])
	if err != nil {
		logrus.Fatalf("Error converting patch version: %s to int\n", splitedPatch[0])
	}
	identifier := strings.Join(splitedPatch[1:], "-")
	if identifier != "" {
		identifier = "-" + identifier
	}
	return major, minor, patch, identifier
}

func getVersionForStableBranch(previousTag string) string {
	if versionMatchBranchName(previousTag, gitService.GetCurrentBranchName()) {
		return bumpPatchVersion(previousTag)
	}
	return getVersionFromBranchName(gitService.GetCurrentBranchName())
}

func getVersionFromBranchName(branch string) string {
	if isStableBranch(branch) {
		regex := regexp.MustCompile("^(.*)-stable")
		match := regex.FindStringSubmatch(branch)
		if len(match) < 2 {
			logrus.Fatalf("Version not found in branch name: %s\n", branch)
		}
		return fmt.Sprintf("%s.0", match[1])
	}
	regex := regexp.MustCompile("(\\d+\\.\\d+\\.\\d+)")
	match := regex.FindStringSubmatch(branch)
	if len(match) < 2 {
		logrus.Fatalf("Version not found in branch name: %s\n", branch)
	}
	return match[1]
}

func versionMatchBranchName(version string, branch string) bool {
	regex := regexp.MustCompile("(\\d+\\.\\d+)")
	match := regex.FindStringSubmatch(branch)
	match2 := regex.FindStringSubmatch(version)
	if len(match) < 2 || len(match2) < 2 {
		logrus.Fatalf("Version not found in branch name: %s\n", branch)
	}
	return match[1] == match2[1]
}

func isStableBranch(branch string) bool {
	match, _ := regexp.MatchString("^.*-stable$", branch)
	match2, _ := regexp.MatchString("^.*-stable$", viper.GetString("CI_COMMIT_REF_NAME"))
	return match || match2
}

func isReleaseOrHotfixBranch(branch string) bool {
	match, _ := regexp.MatchString("hotfix.*$", branch)
	match2, _ := regexp.MatchString("^release.*$", branch)
	return match || match2
}

func stableBranchExists(version string) bool {
	major, minor, _, _ := splitSemver(version)
	return gitService.StableBranchExists(major, minor)
}
