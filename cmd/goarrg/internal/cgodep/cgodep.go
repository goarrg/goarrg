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

package cgodep

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
)

type Dep struct {
	Name           string
	Short          string
	Long           string
	Install        func()
	TargetSpecific bool // whether to use a target specific install dir
}

// flags used by goarrg-config for when go invokes ${PKG_CONFIG}
// goarrg-config will automatically add "-I{{.InstallDir}}/include" and "-L{{.InstallDir}}/lib"
type cgoFlags struct {
	CFlags  string
	LDFlags string

	/*
		StaticLDFlags are the LDFlags that would be passed to cgo if pkg-config
		was executed with the "--static" flag. However, unlike pkg-config,
		goarrg would not combine LDFlags and StaticLDFlags.
	*/
	StaticLDFlags string
}

type cgoDep struct {
	name           string
	short          string
	long           string
	version        string
	targetSpecific bool // whether to use a target specific install dir
	install        func(string) cgoFlags
}

type cgoDepMeta struct {
	Version  string
	CgoFlags cgoFlags
}

const cgoDepMetaFileName = "goarrg_cgodep.json"

var cgoDeps = map[string]cgoDep{}

func cgoDepPath(name string) string {
	return filepath.Join(cmd.ModuleDataPath(), "cgodep", name)
}

func cgoDepCache() string {
	return filepath.Join(cmd.ModuleDataPath(), "cgodep", "cache")
}

func cleanDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(debug.ErrorWrap(err, "Failed to scan directory: %q", dir))
	}

	for _, entry := range entries {
		entry := filepath.Join(dir, entry.Name())
		debug.LogV("Removing: %q", entry)
		if err := os.RemoveAll(entry); err != nil {
			panic(debug.ErrorWrap(err, "Failed to remove: %q", entry))
		}
	}
}

func cleanDirFiltered(dir string, filter func(string, os.DirEntry) bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(debug.ErrorWrap(err, "Failed to scan directory: %q", dir))
	}

	for _, entry := range entries {
		target := filepath.Join(dir, entry.Name())
		if filter(target, entry) {
			debug.LogV("Removing: %q", target)
			if err := os.RemoveAll(target); err != nil {
				panic(debug.ErrorWrap(err, "Failed to remove: %q", target))
			}
		}
	}
}

func install(name string) {
	// should never trigger, but just in case
	if _, ok := cgoDeps[name]; !ok {
		panic(debug.ErrorNew("Unknown cgodep %q", name))
	}

	if !toolchain.CgoSupported() {
		debug.LogE("%q is not available on target: %q: cgo unsupported", name, toolchain.Target())
		os.Exit(2)
	}

	debug.LogI("Installing: %q", name)

	dir := cgoDepPath(name)
	if cgoDeps[name].targetSpecific {
		dir = filepath.Join(dir, toolchain.Target())
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(debug.ErrorWrap(err, "Failed to create directory: %q", dir))
	}

	cleanDir(dir)

	// install and write metadata file
	{
		cgoFlags := cgoDeps[name].install(dir)
		meta := cgoDepMeta{Version: cgoDeps[name].version, CgoFlags: cgoFlags}
		j, err := json.Marshal(meta)
		if err != nil {
			panic(debug.ErrorWrap(err, "json.Marshal failed"))
		}

		metaFile := filepath.Join(dir, cgoDepMetaFileName)
		if err := os.WriteFile(metaFile, j, 0o644); err != nil {
			panic(debug.ErrorWrap(err, "Failed to write: %q", metaFile))
		}
	}

	debug.LogI("Installed: %q", name)
}

func List() []Dep {
	deps := make([]Dep, 0, len(cgoDeps))

	for _, dep := range cgoDeps {
		dep := dep
		deps = append(deps,
			Dep{
				Name:  dep.name,
				Short: dep.short,
				Long:  dep.long,
				Install: func() {
					install(dep.name)
				},

				TargetSpecific: dep.targetSpecific,
			})
	}

	return deps
}

func Check() {
	if !toolchain.CgoSupported() {
		debug.LogI("Cgo unsupported, skipping C dependencies check")
		return
	}

	debug.LogI("Checking installed cgo dependencies")
	installed := false
	for name, d := range cgoDeps {
		debug.LogI("Checking for: %q", name)

		dir := cgoDepPath(name)
		if cgoDeps[name].targetSpecific {
			dir = filepath.Join(dir, toolchain.Target())
		}

		{
			info, err := os.Stat(dir)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					panic(debug.ErrorWrap(err, "Unknown error reading: %q", dir))
				}

				debug.LogI("%q not found", name)
				continue
			}

			if !info.IsDir() {
				panic(debug.ErrorNew("%q is not a directory", dir))
			}
		}
		{
			metaFile := filepath.Join(dir, cgoDepMetaFileName)
			j, err := os.ReadFile(metaFile)

			// only care if error is not os.ErrNotExist, json.Unmarshal will take care of the rest
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				panic(debug.ErrorWrap(err, "Unknown error reading: %q", metaFile))
			}

			meta := cgoDepMeta{}
			if err := json.Unmarshal(j, &meta); err != nil || meta.Version != d.version {
				install(d.name)
				installed = true
			} else {
				debug.LogI("%q is up to date", name)
			}
		}
	}

	if installed {
		// we have to clean the cache everytime we change a external C dependency
		// cause the cache would not detect those changes and so would use the
		// old cached object files.
		if err := exec.Run("go", "clean", "-cache"); err != nil {
			panic(err)
		}
	}
}

func Clean() {
	debug.LogI("Cleaning installed cgo dependencies")
	cleanDirFiltered(filepath.Join(cmd.ModuleDataPath(), "cgodep"), func(target string, entry os.DirEntry) bool {
		if entry.Name() == "cache" {
			return false
		}
		if !entry.IsDir() {
			return true
		}

		debug.LogI("Cleaning: %q", target)
		cleanDirFiltered(target, func(target string, entry os.DirEntry) bool {
			if entry.IsDir() && toolchain.ValidPlatform(entry.Name()) {
				cleanDir(target)
				return false
			}
			return true
		})

		return false
	})
}

func CleanCache() {
	err := os.RemoveAll(cgoDepCache())
	if err != nil {
		panic(debug.ErrorWrap(err, "Failed to clean C dependencies downloaded files"))
	}
}
