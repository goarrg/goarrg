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
	"errors"
	"os"
	"path/filepath"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
)

type ResolveMode uint32

const (
	ResolveExists ResolveMode = 0
	ResolveCFlags ResolveMode = 1 << (iota - 1)
	ResolveLDFlags
	ResolveStaticFlags
)

/*
Resolve searches for the dependencies listed in deps and returns the data requested
by mode.
*/
func Resolve(target toolchain.Target, mode ResolveMode, deps ...string) ([]string, error) {
	searchPath := os.Getenv("CGODEP_PATH")
	searchList := filepath.SplitList(searchPath)
	output := []string{}

	{
		// filter out blank paths
		n := 0
		for _, dir := range searchList {
			if dir != "" {
				searchList[n] = dir
				n++
			}
		}
		searchList = searchList[:n]
	}

	if len(searchList) == 0 {
		return nil, debug.Errorf("CGODEP_PATH unset/empty")
	}

	{
		// pick out unique deps
		depMap := map[string]struct{}{}
		for _, d := range deps {
			depMap[d] = struct{}{}
		}

		// filter out unique deps keeping order
		n := 0
		for _, d := range deps {
			if _, ok := depMap[d]; ok {
				deps[n] = d
				delete(depMap, d)
				n++
			}
		}
		deps = deps[:n]
	}

	for _, d := range deps {
		var m Meta
		var err error

		for _, search := range searchList {
			dir := filepath.Join(search, d, target.String())
			m, err = ReadMetaFile(dir)
			if err == nil {
				break
			}
			if !errors.Is(err, os.ErrNotExist) { // dep found but failed to load
				return nil, debug.ErrorWrapf(err, "Failed to resolve %q", d)
			}

			dir = filepath.Join(search, d)
			m, err = ReadMetaFile(dir)
			if err == nil {
				break
			}
			if !errors.Is(err, os.ErrNotExist) { // dep found but failed to load
				return nil, debug.ErrorWrapf(err, "Failed to resolve %q", d)
			}
		}

		if err != nil { // dep not found anywhere
			return nil, debug.ErrorWrapf(debug.Errorf("Dependency not installed"), "Failed to resolve %q", d)
		}

		if (mode & ResolveCFlags) == ResolveCFlags {
			output = append(output, m.Flags.CFlags...)
		}
		if (mode & ResolveLDFlags) == ResolveLDFlags {
			if (mode & ResolveStaticFlags) == ResolveStaticFlags {
				output = append(output, m.Flags.StaticLDFlags...)
			} else {
				output = append(output, m.Flags.LDFlags...)
			}
		}
	}

	return output, nil
}
