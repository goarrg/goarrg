/*
Copyright 2021 The goARRG Authors.

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
	"path/filepath"

	"goarrg.com/debug"
)

func cleanDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to scan directory: %q", dir))
	}

	for _, entry := range entries {
		entry := filepath.Join(dir, entry.Name())
		debug.VPrintf("Removing: %q", entry)
		if err := os.RemoveAll(entry); err != nil {
			panic(debug.ErrorWrapf(err, "Failed to remove: %q", entry))
		}
	}
}

/*
ToolsDir returns a path within WorkingModuleDataDir() to be used to store tools
needed to build the module.
*/
func ToolsDir() string {
	return filepath.Join(WorkingModuleDataDir(), "tool")
}

/*
CleanTools clears out the folder located at ToolsDir()
*/
func CleanTools() {
	debug.IPrintf("Cleaning tools folder")
	cleanDir(ToolsDir())
}
