/*
Copyright 2022 The goARRG Authors.

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
	"path/filepath"

	"goarrg.com/debug"
	"golang.org/x/tools/go/packages"
)

var (
	modulePath    string
	moduleDataDir string
)

func init() {
	p, err := packages.Load(&packages.Config{Mode: packages.NeedModule}, ".")
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to load package in current directory"))
	}
	if len(p) == 0 {
		panic(debug.Errorf("No go package in current directory"))
	}

	// if len is > 1 the module should still all be the same
	modulePath = p[0].Module.Path
	moduleDataDir = filepath.Join(p[0].Module.Dir, ".goarrg")
}

/*
WorkingModulePath returns the module path of the module in the working
directory at the time of init.
*/
func WorkingModulePath() string {
	return modulePath
}

/*
WorkingModuleDataDir returns a folder within WorkingModulePath() that can be
used to store data.
*/
func WorkingModuleDataDir() string {
	return moduleDataDir
}
