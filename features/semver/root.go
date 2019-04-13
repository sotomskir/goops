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

const (
	// Configuration variables
	GoopscSemver           = "GOOPSC_SEMVER"
	GoopscSemverStrategy   = "GOOPSC_SEMVER_STRATEGY"
	GoopscSemverSaveExport = "GOOPSC_SAVE_EXPORT"

	// Output variables
	GoopsSemver           = "GOOPS_SEMVER"
	GoopsSemverRelease    = "GOOPS_SEMVER_RELEASE"
	GoopsSemverMajor      = "GOOPS_SEMVER_MAJOR"
	GoopsSemverMinor      = "GOOPS_SEMVER_MINOR"
	GoopsSemverPatch      = "GOOPS_SEMVER_PATCH"

	// Configuration options
	GithubFlowStrategy    = "github-flow"
	GitlabFlowStrategy    = "gitlab-flow"
	GitFlowBranchStrategy = "git-flow-branch"
)

func setDefaults() {
	viper.SetDefault(GoopscSemver, "false")
	viper.SetDefault(GoopscSemverSaveExport, "true")
	viper.SetDefault(GoopscSemverStrategy, GithubFlowStrategy)
}

type strategy interface {
	getSemanticVersion() string
}

type Semver struct {
	strategy strategy
}

func New() Semver {
	setDefaults()
	var strategy strategy
	switch viper.GetString(GoopscSemverStrategy) {
	case GithubFlowStrategy:
		strategy = githubFlow{}
	case GitlabFlowStrategy:
		strategy = gitlabFlow{}
	case GitFlowBranchStrategy:
		strategy = gitFlowBranch{}
	default:
		logrus.Errorf("Unexpected strategy: %s\n", viper.GetString(GoopscSemverStrategy))
		os.Exit(1)
	}
	return Semver{strategy: strategy}
}

func (o *Semver) GetVersion() string {
	if utils.IsDisabled(GoopscSemver) {
		return ""
	}
	version := o.strategy.getSemanticVersion()
	major, minor, patch, _ := splitSemver(version)
	if utils.IsEnabled(GoopscSemverSaveExport) {
		utils.SaveExportString(GoopsSemver, version)
		utils.SaveExportString(GoopsSemverRelease, strings.Replace(version, "-SNAPSHOT", "", 1))
		utils.SaveExportInt(GoopsSemverMajor, major)
		utils.SaveExportInt(GoopsSemverMinor, minor)
		utils.SaveExportInt(GoopsSemverPatch, patch)
	}
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
		logrus.Fatalf("Error converting minor version: %s to int\n", splited[1])
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
		match := findMatch(regex, branch)
		return fmt.Sprintf("%s.0", match[1])
	}
	regex := regexp.MustCompile("(\\d+\\.\\d+\\.\\d+)")
	match := findMatch(regex, branch)
	return match[1]
}

func findMatch(regex *regexp.Regexp, branch string) []string {
	match := regex.FindStringSubmatch(branch)
	if len(match) < 2 {
		logrus.Fatalf("Version not found in branch name: %s\n", branch)
	}
	return match
}

func versionMatchBranchName(version string, branch string) bool {
	regex := regexp.MustCompile("(\\d+\\.\\d+)")
	match := findMatch(regex, branch)
	match2 := findMatch(regex, version)
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
