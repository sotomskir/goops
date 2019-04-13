// Copyright © 2019 Robert Sotomski <sotomski@gmail.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"github.com/sotomskir/goops/features/jira"
	"github.com/sotomskir/goops/features/semver"
	"github.com/spf13/cobra"
)

var (
	summary string
	description string
	issueType string
)

// setenvCmd represents the pipelineCommon command
var setenvCmd = &cobra.Command{
	Use:     "setenv",
	Aliases: []string{"s"},
	Short:   "Sets environment variables, and runs common tasks. Should be called prior to other commands",
	Long: `Sets environment variables, and runs common tasks. Should be called prior to other commands.
This command will call version and issues commands internally. And will save following variables to gitlab.env file:
CI_SEMVER_RELEASE, CI_SEMVER, CI_ISSUES. This way allows to share variables between pipeline stages.
Variables can be used in next stages by reading them from gitlab.env file, using command "source .goops.env". 
If build is not in merge context CI_ISSUES will be fetched from previous merged merge request.
All CI_ISSUES will be assigned to CI_SEMVER_RELEASE version in Jira. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		s := semver.New()
		j := jira.New()
		version := s.GetVersion()
		issues := j.GetIssues()
		j.SetJiraVersion(version, issues, summary, description, issueType)
	},
}

func init() {
	rootCmd.AddCommand(setenvCmd)
	rootCmd.Flags().StringVarP(&summary, "summary", "s", "", "Deployment issue summary.")
	rootCmd.Flags().StringVarP(&description, "description", "d", "", "Deployment issue description.")
	rootCmd.Flags().StringVarP(&issueType, "issue-type", "t", "", "Deployment issue type.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Pe rsistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setenvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setenvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
