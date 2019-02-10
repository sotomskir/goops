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

package execService

import (
	"os/exec"
	"strings"
)

type IService interface {
	Exec(cmd string) (string, error)
}

type Service struct {
}

func (s Service) Exec(cmd string) (string, error) {
	args := strings.Split(cmd, " ")
	name := args[0]
	arg := args[1:]
	out, err := exec.Command(name, arg...).CombinedOutput()
	return strings.Trim(string(out), " \n"), err
}
