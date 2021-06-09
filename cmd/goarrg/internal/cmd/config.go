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

package cmd

import (
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

var (
	cwd             string
	moduleDataPath  string
	packageMain     bool
	flagVerbose     bool
	flagVeryVerbose bool
)

func init() {
	var err error
	cwd, err = os.Getwd()

	if err != nil {
		panic(err)
	}

	p, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedModule}, ".")
	if err != nil {
		panic(err)
	}

	if len(p) < 1 {
		panic("No go package in current directory")
	}

	if len(p) > 1 {
		panic("Multiple go packages in current directory")
	}

	packageMain = p[0].Name == "main"
	moduleDataPath = filepath.Join(p[0].Module.Dir, ".goarrg")
}

func CWD() string {
	return cwd
}

func ResetCWD() {
	if err := os.Chdir(cwd); err != nil {
		panic(err)
	}
}

func ModuleDataPath() string {
	return moduleDataPath
}

func PackageMain() bool {
	return packageMain
}

func Verbose() bool {
	return flagVerbose || flagVeryVerbose
}

func VeryVerbose() bool {
	return flagVeryVerbose
}
