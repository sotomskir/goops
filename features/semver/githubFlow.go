package semver

import (
	"fmt"
	"github.com/sotomskir/goops/gitService"
)

type githubFlow struct{}

func (githubFlow) getSemanticVersion() string {
	headTag := gitService.GetHeadTag()
	if headTag != "" {
		return headTag
	}
	previousTag := gitService.GetPreviousTag()
	var version string
	if previousTag == "" {
		previousTag = "0.0.0"
	}
	version = bumpMinorVersion(previousTag)
	return fmt.Sprintf("%s-SNAPSHOT", version)
}
