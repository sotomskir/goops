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
	"github.com/sotomskir/goops/features/docker"
	"github.com/spf13/cobra"
)

// pipelineDockerBuildCmd represents the pipelineDockerBuild command
var pipelineDockerBuildCmd = &cobra.Command{
	Use:     "build PATH",
	Aliases: []string{"b"},
	Short:   "Build docker image",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dockerfile, _ := cmd.Flags().GetString("file")
		tag, err := cmd.Flags().GetString("tag")
		if err != nil {
			logrus.Fatalln(err)
		}
		docker.DockerBuild(tag, dockerfile, args[0])
	},
}

func init() {
	pipelineDockerCmd.AddCommand(pipelineDockerBuildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineDockerBuildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pipelineDockerBuildCmd.Flags().StringP("file", "f", "Dockerfile", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	pipelineDockerBuildCmd.Flags().StringP("tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	pipelineDockerBuildCmd.MarkFlagRequired("tag")
}
