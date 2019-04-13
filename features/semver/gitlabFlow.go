package semver

import (
	"fmt"
	"github.com/sotomskir/goops/gitService"
)

type gitlabFlow struct{}

func (gitlabFlow) getSemanticVersion() string {
	headTag := gitService.GetHeadTag()
	if headTag != "" {
		return headTag
	}
	previousTag := gitService.GetPreviousTag()
	var version string
	if previousTag == "" {
		previousTag = "0.0.0"
	}
	if isStableBranch(gitService.GetCurrentBranchName()) {
		version = getVersionForStableBranch(previousTag)
	} else {
		version = bumpMinorVersion(previousTag)
		if stableBranchExists(version) {
			version = bumpMinorVersion(version)
		}
	}
	return fmt.Sprintf("%s-SNAPSHOT", version)
}
