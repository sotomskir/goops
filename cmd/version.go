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
	"github.com/sotomskir/gitlab-cli/gitService"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Aliases: []string{"v"},
	Short: "Generate semantic version for current HEAD",
	Long: `Generate semantic version for current HEAD.
Version generation is based on git tags.
If current HEAD is tagged then tag will be used as version.
Else command will lookup for previous tag bump it's minor version, reset patch version and append '-SNAPSHOT'
When there are no tags found version will be '0.1.0-SNAPSHOT'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(gitService.GetSemanticVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
