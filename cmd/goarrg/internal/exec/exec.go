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
	"path/filepath"
	"runtime"
	"strings"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/debug"
)

func Run(c string, args ...string) error {
	{
		cmd := c
		for _, arg := range args {
			if strings.Contains(arg, " ") {
				cmd += " \"" + arg + "\""
			} else {
				cmd += " " + arg
			}
		}
		debug.LogI("Running: %s", cmd)
	}

	cmd := exec.Command(c, args...)

	if runtime.GOOS == "windows" && c == "cmake" {
		path := os.Getenv("PATH")
		specialPath := ""

		for _, p := range filepath.SplitList(path) {
			if _, err := os.Stat(filepath.Join(p, "sh.exe")); os.IsNotExist(err) {
				specialPath = specialPath + string(filepath.ListSeparator) + p
			}
		}

		os.Setenv("PATH", specialPath[1:])
		cmd.Env = os.Environ()
		os.Setenv("PATH", path)
	} else {
		cmd.Env = os.Environ()
	}

	if base.IsVeryVerbose() {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr

	return debug.ErrorWrap(cmd.Run(), "Failed to run command")
}

func RunOutput(c string, args ...string) ([]byte, error) {
	{
		cmd := c
		for _, arg := range args {
			if strings.Contains(arg, " ") {
				cmd += " \"" + arg + "\""
			} else {
				cmd += " " + arg
			}
		}
		debug.LogI("Running: %s", cmd)
	}

	cmd := exec.Command(c, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	return out, debug.ErrorWrap(err, "Failed to run command")
}

func RunExit(c string, args ...string) int {
	{
		cmd := c
		for _, arg := range args {
			if strings.Contains(arg, " ") {
				cmd += " \"" + arg + "\""
			} else {
				cmd += " " + arg
			}
		}
		debug.LogI("Running: %s", cmd)
	}

	cmd := exec.Command(c, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr

	if base.IsVeryVerbose() {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Run()

	if err == nil {
		return 0
	}

	return err.(*exec.ExitError).ExitCode()
}
