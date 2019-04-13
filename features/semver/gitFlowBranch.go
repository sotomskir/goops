package semver

import (
	"fmt"
	"github.com/sotomskir/goops/gitService"
)

type gitFlowBranch struct{}

func (gitFlowBranch) getSemanticVersion() string {
	previousMergedVersion, err := gitService.GetPreviouslyMergedVersion()
	if err != nil {
		panic(err)
	}
	branch := gitService.GetCurrentBranchName()
	if branch == "master" {
		return previousMergedVersion
	}
	var version string
	if isReleaseOrHotfixBranch(branch) {
		version = getVersionFromBranchName(branch)
	} else {
		version = bumpMinorVersion(previousMergedVersion)
		if gitService.BranchExists(version) {
			version = bumpMinorVersion(version)
		}
	}
	return fmt.Sprintf("%s-SNAPSHOT", version)
}
