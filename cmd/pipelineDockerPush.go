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
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/gitlab-cli/pipelineApi"

	"github.com/spf13/cobra"
)

// pipelineDockerPushCmd represents the pipelineDockerPush command
var pipelineDockerPushCmd = &cobra.Command{
	Use:   "push TAG",
	Aliases: []string{"p"},
	Short: "Push docker images to registry",
	Long: `Push docker images to registry. 
If build context is not one of: master, tags, ^.*-stable$ push will be skipped.
If build is from git tag it will also push image with "stable" tag.
If build is from master branch it will also push image with "latest" tag`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tag, err := cmd.Flags().GetString("tag")
		if err != nil {
			logrus.Fatalln(err)
		}
		pipelineApi.DockerPush(tag)
	},
}

func init() {
	pipelineDockerCmd.AddCommand(pipelineDockerPushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineDockerPushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pipelineDockerPushCmd.Flags().StringP("tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
}
