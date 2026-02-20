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
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"goarrg.com/build"
	"goarrg.com/debug"
)

type Build = build.Build

const (
	BuildRelease     = build.BuildRelease
	BuildDevelopment = build.BuildDevelopment
	BuildDebug       = build.BuildDebug
)

type Target struct {
	OS   string
	Arch string
}

func (t Target) String() string {
	if (t == Target{}) {
		return runtime.GOOS + "_" + runtime.GOARCH
	}
	return t.OS + "_" + t.Arch
}

func (t Target) CrossCompiling() bool {
	return (t != Target{}) && (t.OS != runtime.GOOS || t.Arch != runtime.GOARCH)
}

/*
IgnoreBlacklist returns an ignore function that returns true if the path is one of the args passed to IgnoreBlacklist.
*/
func IgnoreBlacklist(args ...string) func(string) bool {
	blacklist := map[string]struct{}{}
	for _, arg := range args {
		blacklist[arg] = struct{}{}
	}
	return func(path string) bool {
		_, skip := blacklist[path]
		return skip
	}
}

/*
ScanDirModTime scans recursively dir and returns the last time a file was modified or error.
ignore is a function that takes in a path relative to dir and if it returns true, the path is skipped.
If dir does not exist, returns time.Unix(0, 0)
*/
func ScanDirModTime(dir string, ignore func(string) bool) time.Time {
	latestMod := time.Unix(0, 0)
	err := filepath.Walk(dir, func(path string, fs fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == dir {
			return err
		}
		{
			rel := strings.TrimPrefix(path, dir+string(filepath.Separator))
			if ignore != nil && ignore(rel) {
				if fs.IsDir() {
					return filepath.SkipDir
				}
				return err
			}
		}
		if fs.IsDir() {
			return err
		}
		mod := fs.ModTime()
		if mod.After(latestMod) {
			latestMod = mod
		}
		return err
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(debug.ErrorWrapf(err, "Failed to scan %q", dir))
	}
	return latestMod
}
