/*
Copyright 2020 The goARRG Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exec

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/debug"
)

func Run(c string, args ...string) error {
	{
		cmd := c
		for _, arg := range args {
			if strings.Contains(arg, " ") {
				cmd += " " + strconv.Quote(arg)
			} else {
				cmd += " " + arg
			}
		}
		debug.LogI("Running: %s", cmd)
	}

	ex := exec.Command(c, args...)
	ex.Env = os.Environ()
	ex.Stderr = os.Stderr

	if cmd.VeryVerbose() {
		ex.Stdout = os.Stdout
	}

	return debug.ErrorWrap(ex.Run(), "Failed to run command")
}

func RunOutput(c string, args ...string) ([]byte, error) {
	{
		cmd := c
		for _, arg := range args {
			if strings.Contains(arg, " ") {
				cmd += " " + strconv.Quote(arg)
			} else {
				cmd += " " + arg
			}
		}
		debug.LogI("Running: %s", cmd)
	}

	ex := exec.Command(c, args...)
	ex.Env = os.Environ()
	ex.Stderr = os.Stderr

	out, err := ex.Output()
	return out, debug.ErrorWrap(err, "Failed to run command")
}

func RunExit(c string, args ...string) int {
	{
		cmd := c
		for _, arg := range args {
			if strings.Contains(arg, " ") {
				cmd += " " + strconv.Quote(arg)
			} else {
				cmd += " " + arg
			}
		}
		debug.LogI("Running: %s", cmd)
	}

	ex := exec.Command(c, args...)
	ex.Env = os.Environ()
	ex.Stderr = os.Stderr

	if cmd.VeryVerbose() {
		ex.Stdout = os.Stdout
	}

	err := ex.Run()

	if err == nil {
		return 0
	}

	return err.(*exec.ExitError).ExitCode()
}

func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
