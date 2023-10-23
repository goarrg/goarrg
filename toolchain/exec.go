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

package toolchain

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"goarrg.com/debug"
)

func printRun(c string, args []string) {
	cmd := c
	for _, arg := range args {
		if strings.Contains(arg, " ") {
			cmd += " " + strconv.Quote(arg)
		} else {
			cmd += " " + arg
		}
	}
	debug.IPrintf("Running: %s", cmd)
}

func RunEnvDir(env []string, dir, c string, args ...string) error {
	printRun(c, args)
	ex := exec.Command(c, args...)
	ex.Env = env
	ex.Dir = dir
	ex.Stderr = os.Stderr

	if debug.WillLog(debug.LogLevelVerbose) {
		ex.Stdout = os.Stdout
	}

	return debug.ErrorWrapf(ex.Run(), "Failed to run command")
}

func RunEnvDirCombinedOutput(env []string, dir, c string, args ...string) ([]byte, error) {
	printRun(c, args)
	ex := exec.Command(c, args...)
	ex.Env = env
	ex.Dir = dir

	out, err := ex.CombinedOutput()
	return out, debug.ErrorWrapf(err, "Failed to run command")
}

func RunEnvDirOutput(env []string, dir, c string, args ...string) ([]byte, error) {
	printRun(c, args)
	ex := exec.Command(c, args...)
	ex.Env = env
	ex.Dir = dir
	ex.Stderr = os.Stderr

	out, err := ex.Output()
	return out, debug.ErrorWrapf(err, "Failed to run command")
}

func RunEnvDirExit(env []string, dir, c string, args ...string) int {
	printRun(c, args)
	ex := exec.Command(c, args...)
	ex.Env = env
	ex.Dir = dir
	ex.Stderr = os.Stderr

	if debug.WillLog(debug.LogLevelVerbose) {
		ex.Stdout = os.Stdout
	}

	err := ex.Run()

	if err == nil {
		return 0
	}

	return err.(*exec.ExitError).ExitCode() // nolint: errorlint
}

func RunDir(dir, c string, args ...string) error {
	return RunEnvDir(os.Environ(), dir, c, args...)
}

func RunDirCombinedOutput(dir, c string, args ...string) ([]byte, error) {
	return RunEnvDirCombinedOutput(os.Environ(), dir, c, args...)
}

func RunDirOutput(dir, c string, args ...string) ([]byte, error) {
	return RunEnvDirOutput(os.Environ(), dir, c, args...)
}

func RunDirExit(dir, c string, args ...string) int {
	return RunEnvDirExit(os.Environ(), dir, c, args...)
}

func RunEnv(env []string, c string, args ...string) error {
	return RunEnvDir(env, "", c, args...)
}

func RunEnvCombinedOutput(env []string, c string, args ...string) ([]byte, error) {
	return RunEnvDirCombinedOutput(env, "", c, args...)
}

func RunEnvOutput(env []string, c string, args ...string) ([]byte, error) {
	return RunEnvDirOutput(env, "", c, args...)
}

func RunEnvExit(env []string, c string, args ...string) int {
	return RunEnvDirExit(env, "", c, args...)
}

func Run(c string, args ...string) error {
	return RunEnvDir(os.Environ(), "", c, args...)
}

func RunCombinedOutput(c string, args ...string) ([]byte, error) {
	return RunEnvDirCombinedOutput(os.Environ(), "", c, args...)
}

func RunOutput(c string, args ...string) ([]byte, error) {
	return RunEnvDirOutput(os.Environ(), "", c, args...)
}

func RunExit(c string, args ...string) int {
	return RunEnvDirExit(os.Environ(), "", c, args...)
}
