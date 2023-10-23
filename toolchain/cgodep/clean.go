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

package cgodep

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

func cleanDirFiltered(dir string, filter func(string, os.DirEntry) bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to scan directory: %q", dir))
	}

	for _, entry := range entries {
		target := filepath.Join(dir, entry.Name())
		if filter(target, entry) {
			debug.VPrintf("Removing: %q", target)
			if err := os.RemoveAll(target); err != nil {
				panic(debug.ErrorWrapf(err, "Failed to remove: %q", target))
			}
		}
	}
}

/*
CleanCache clears out the download cache but not built dependencies.
*/
func CleanCache() {
	debug.IPrintf("Cleaning cgodep cache folder")
	cleanDir(cacheDir())
}

/*
Clean clears out the built dependencies but not the download cache.
*/
func Clean() {
	debug.IPrintf("Cleaning installed cgo dependencies")
	cleanDirFiltered(DataDir(), func(target string, entry os.DirEntry) bool {
		return entry.Name() != "cache"
	})
}
