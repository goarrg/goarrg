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
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"goarrg.com/debug"
	"golang.org/x/tools/go/packages"
)

var (
	modulePath    string
	moduleDir     string
	moduleDataDir string
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to get current directory"))
	}
	SetWorkingModule(cwd)
}

/*
SetWorkingModule sets the working module to what module dir belongs to if found.
If no module is found then module dir will be set to dir with WorkingModulePath() an empty string.
It only affects the WorkingModule[Path|Dir|DataDir] functions.
It does not actually change os.Getwd().
*/
func SetWorkingModule(dir string) {
	// packages[0].Module will be null if dir has no go files so we need to search child dirs too
	p, err := packages.Load(&packages.Config{Mode: packages.NeedModule}, filepath.Join(dir, "..."))
	if err != nil && len(p) > 0 {
		m := map[string]*packages.Module{}
		for _, pkg := range p {
			if pkg.Module != nil {
				// there may be multiple modules, we store the relative path
				// so we can later pick the nearest one
				if r, err := filepath.Rel(pkg.Module.Dir, dir); err == nil {
					m[r] = pkg.Module
				}
			}
		}

		if len(m) > 0 {
			// sort based on how deep in the filetree we are
			module := m[slices.SortedFunc(maps.Keys(m), func(a, b string) int {
				return strings.Count(filepath.ToSlash(a), "/") - strings.Count(filepath.ToSlash(b), "/")
			})[0]]
			modulePath = module.Path
			moduleDir = module.Dir
			moduleDataDir = filepath.Join(moduleDir, ".goarrg")
			return
		}
	}

	// we are not in a module
	modulePath = ""
	moduleDir = dir
	moduleDataDir = filepath.Join(moduleDir, ".goarrg")
}

/*
WorkingModulePath returns the module path of the module in the last call to
SetWorkingModule or an empty string if no module is found.
*/
func WorkingModulePath() string {
	return modulePath
}

/*
WorkingModulePath returns the filepath of the module in the last call to
SetWorkingModule, if no module was found this will be equal to what was
passed to SetWorkingModule.
*/
func WorkingModuleDir() string {
	return moduleDir
}

/*
WorkingModuleDataDir returns a folder within WorkingModuleDir() that can be
used to store data.
*/
func WorkingModuleDataDir() string {
	return moduleDataDir
}
