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

package base

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var envGOOS, envGOARCH string
var usrCache, usrData string
var cwd string

var flagVerbose bool
var flagVeryVerbose bool

func init() {
	if envGO := os.Getenv("GO"); envGO != "" {
		cmd := exec.Command(envGO, "env", "GOROOT")
		cmd.Env = os.Environ()

		out, err := cmd.Output()

		if err != nil {
			panic(err)
		}

		goroot := strings.TrimSpace(string(out))
		path := os.Getenv("PATH")

		if err := os.Setenv("PATH", filepath.Join(goroot, "bin")+string(filepath.ListSeparator)+path); err != nil {
			panic(err)
		}
	}

	envGOOS = os.Getenv("GOOS")
	envGOARCH = os.Getenv("GOARCH")

	if envGOOS == "" {
		envGOOS = runtime.GOOS
		os.Setenv("GOOS", envGOOS)
	}

	if envGOARCH == "" {
		envGOARCH = runtime.GOARCH
		os.Setenv("GOARCH", envGOARCH)
	}

	var err error
	usrCache, err = os.UserCacheDir()

	if err != nil {
		panic(err)
	}

	usrCache = filepath.Join(usrCache, "goarrg")

	switch runtime.GOOS {
	case "linux":
		usrData = os.Getenv("XDG_DATA_HOME")
		if usrData == "" {
			usrHome, err := os.UserHomeDir()

			if err != nil {
				panic(err)
			}

			usrData = filepath.Join(usrHome, ".local", "share", "goarrg")
		} else {
			usrData = filepath.Join(usrData, "goarrg")
		}
	case "windows":
		usrData = usrCache
		usrCache = filepath.Join(usrCache, "cache")
	default:
		panic("No support for " + runtime.GOOS)
	}

	cwd, err = os.Getwd()

	if err != nil {
		panic(err)
	}
}

func GOARCH() string {
	return envGOARCH
}

func GOOS() string {
	return envGOOS
}

func UsrCache() string {
	return usrCache
}

func UsrData() string {
	return usrData
}

func CWD() string {
	return cwd
}

func ResetCWD() {
	if err := os.Chdir(cwd); err != nil {
		panic(err)
	}
}

func IsVerbose() bool {
	return flagVerbose
}

func IsVeryVerbose() bool {
	return flagVeryVerbose
}
