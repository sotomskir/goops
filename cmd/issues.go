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
	"fmt"
	"github.com/sotomskir/gitlab-cli/gitlabApi"
	"github.com/sotomskir/gitlab-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issues",
	Short: "List Jira issue keys mentioned in merge request title, description and commit messages",
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		utils.ViperValidate("ci_merge_request_iid", "mr", "CI_MERGE_REQUEST_IID")
		utils.ViperValidate("ci_project_id", "project", "CI_PROJECT_ID")
		issueKeys := gitlabApi.GetMergeRequestIssueKeys(viper.GetString("ci_project_id"), viper.GetString("ci_merge_request_iid"))
		for _, v := range issueKeys {
			fmt.Println(v)
		}
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	issueCmd.Flags().StringP( "project", "p", "", "Project id")
	issueCmd.Flags().StringP( "mr", "m", "", "Merge request iid")

	viper.BindPFlag("ci_project_id", issueCmd.Flags().Lookup("project"))
	viper.BindPFlag("ci_merge_request_iid", issueCmd.Flags().Lookup("mr"))
}
