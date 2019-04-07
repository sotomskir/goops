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
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var e IService = Service{}

type IService interface {
	Exec(cmd string) (string, error)
	LogExec(cmd string)
}

type Service struct {
}

func Exec(cmd string) (string, error) {
	return e.Exec(cmd)
}

func (s Service) Exec(cmd string) (string, error) {
	args := strings.Split(cmd, " ")
	name := args[0]
	arg := args[1:]
	out, err := exec.Command(name, arg...).CombinedOutput()
	return strings.Trim(string(out), " \n"), err
}

func (s Service) LogExec(command string) {
	args := strings.Split(command, " ")
	name := args[0]
	arg := args[1:]
	fmt.Println(command)
	cmd := exec.Command(name, arg...)
	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	outScanner := bufio.NewScanner(cmdOutReader)
	go func() {
		for outScanner.Scan() {
			fmt.Println(outScanner.Text())
		}
	}()
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Println(errScanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func LogExec(cmd string) {
	e.LogExec(cmd)
}
