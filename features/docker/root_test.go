// Copyright © 2019 Robert Sotomski <sotomski@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker

import (
	"github.com/golang/mock/gomock"
	"github.com/sotomskir/goops/mockExecService"
	"github.com/spf13/viper"
	"testing"
)

func TestGetHeadTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIService := mock_execService.NewMockIService(ctrl)
	Initialize(mockIService)

	// If build context is not one of: master, tags, ^.*-stable$ push will be skipped.
	viper.Set("CI_COMMIT_REF_NAME", "feature-test")
	viper.Set("CI_COMMIT_TAG", "")
	DockerPush("test/test:1.0.0")

	// If build is from git tag it will also push image with "stable" tag.
	mockIService.EXPECT().LogExec(gomock.Eq("docker tag test/test:2.0.0 test/test:stable")).Times(1)
	mockIService.EXPECT().LogExec(gomock.Eq("docker push test/test:stable")).Times(1)
	mockIService.EXPECT().LogExec(gomock.Eq("docker push test/test:2.0.0")).Times(1)
	viper.Set("CI_COMMIT_REF_NAME", "2.0-stable")
	viper.Set("CI_COMMIT_TAG", "2.0.0")
	DockerPush("test/test:2.0.0")

	// If build is from master branch it will also push image with "latest" tag`,
	mockIService.EXPECT().LogExec(gomock.Eq("docker tag test/test:3.0.0-SNAPSHOT test/test:latest")).Times(1)
	mockIService.EXPECT().LogExec(gomock.Eq("docker push test/test:latest")).Times(1)
	mockIService.EXPECT().LogExec(gomock.Eq("docker push test/test:3.0.0-SNAPSHOT")).Times(1)
	viper.Set("CI_COMMIT_REF_NAME", "master")
	viper.Set("CI_COMMIT_TAG", "")
	DockerPush("test/test:3.0.0-SNAPSHOT")

	mockIService.EXPECT().LogExec(gomock.Eq("docker push test/test:4.0.1-SNAPSHOT")).Times(1)
	viper.Set("CI_COMMIT_REF_NAME", "4.0-stable")
	viper.Set("CI_COMMIT_TAG", "")
	DockerPush("test/test:4.0.1-SNAPSHOT")
}
