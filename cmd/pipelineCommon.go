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
	"github.com/sotomskir/gitlab-cli/pipelineApi"

	"github.com/spf13/cobra"
)

// pipelineCommonCmd represents the pipelineCommon command
var pipelineCommonCmd = &cobra.Command{
	Use:   "common",
	Aliases: []string{"c"},
	Short: "Sets environment variables, and runs common tasks. Should be called prior to other pipeline commands",
	Long: `Sets environment variables, and runs common tasks. Should be called prior to other pipeline commands.
This command will call version and issues commands internally. And will save following variables to gitlab.env file:
CI_SEMVER_RELEASE, CI_SEMVER, CI_ISSUES. This way allows to share variables between pipeline stages.
Variables can be used in next stages by reading them from gitlab.env file, using command "source gitlab.env". 
If build is not in merge context CI_ISSUES will be fetched from previous merged merge request.
All CI_ISSUES will be assigned to CI_SEMVER_RELEASE version in Jira. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		pipelineApi.Common()
	},
}

func init() {
	pipelineCmd.AddCommand(pipelineCommonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineCommonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipelineCommonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
