/*
Copyright 2026 The goARRG Authors.

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

package cc

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
)

type BuildType uint

const (
	BuildTypeStaticLibrary = iota
	BuildTypeSharedLibrary
	BuildTypeExecuteable
)

func (t BuildType) FileExt(os string) string {
	switch t {
	case BuildTypeStaticLibrary:
		return ".a"
	case BuildTypeSharedLibrary:
		if os == "windows" {
			return ".dll"
		}
		return ".so"
	case BuildTypeExecuteable:
		if os == "windows" {
			return ".exe"
		}
		return ""
	}
	// should never be here
	panic(fmt.Sprintf("Unknown build type: %d", t))
}

type BuildFlags struct {
	CFlags   []string
	CXXFlags []string
}

type BuildOptions struct {
	Type   BuildType
	Target toolchain.Target
	// Flags that are both passed when compiling and to the returned CompileCommand
	BuildFlags BuildFlags
	// Flags that are only passed when compiling, will appear before BuildFlags, will not appear in CompileCommand
	CompileOnlyFlags BuildFlags
	// Flags that are only passed to the returned CompileCommand, will appear before BuildFlags
	CommandOnlyFlags BuildFlags
	// Flags that are only passed when linking
	LDFlags []string

	// if Ignore(path) is true, will skip over dir or file when compiling,
	// path will be a path relative to srcDir
	Ignore func(string) bool
}

type CompileCommand struct {
	Directory string   `json:"directory"`
	Arguments []string `json:"arguments"`
	File      string   `json:"file"`
}

/*
BuildDir scans srcDir for .c and .cpp files to compile with ${CC} or ${CXX} respectively.
If building a static library, ${AR} will be used to generate outFile.
Else linking will be done with ${CC} if no .cpp files were found, otherwise ${CXX}.
Returns an array of CompileCommands if successful, it is meant to be used to generate compile_commands.json for linters.
BuildDir does not generate it automatically as one might want to call it multiple times with varying BuildOptions.Type
and BuildOptions.Ignore functions, generating automatically would wipe past builds.
*/
func BuildDir(srcDir, buildDir, outFile string, options BuildOptions) ([]CompileCommand, error) {
	if options.Target == (toolchain.Target{}) {
		options.Target = toolchain.Target{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		}
	}
	if options.Type == BuildTypeSharedLibrary && options.Target.OS == "linux" {
		// clip to force realloc cause we do not want to risk modifying the original allocation
		options.BuildFlags.CFlags = append(slices.Clip(options.BuildFlags.CFlags), "-fPIC")
		options.BuildFlags.CXXFlags = append(slices.Clip(options.BuildFlags.CXXFlags), "-fPIC")
		options.LDFlags = append(slices.Clip(options.LDFlags), "-fPIC")
	}
	var jsonCmds []CompileCommand
	var compileCmds []CompileCommand
	var objs []string
	var isCXX bool
	{
		err := filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			path = strings.TrimPrefix(path, srcDir+string(filepath.Separator))
			if options.Ignore != nil && options.Ignore(path) {
				return err
			}
			switch filepath.Ext(path) {
			case ".c":
				objs = append(objs, path+".o")
				compileArgs := append(slices.Clip(options.CompileOnlyFlags.CFlags), options.BuildFlags.CFlags...)
				compileCmds = append(compileCmds, CompileCommand{Directory: srcDir, Arguments: append([]string{os.Getenv("CC")},
					compileArgs...), File: path})
				fallthrough
			case ".h":
				jsonArgs := append(slices.Clip(options.CommandOnlyFlags.CFlags), options.BuildFlags.CFlags...)
				jsonArgs = append(jsonArgs, "-c", path)
				jsonCmds = append(jsonCmds,
					CompileCommand{Directory: srcDir, Arguments: append([]string{os.Getenv("CC")},
						jsonArgs...), File: path})
			case ".cpp":
				isCXX = true
				objs = append(objs, path+".o")
				compileArgs := append(slices.Clip(options.CompileOnlyFlags.CXXFlags), options.BuildFlags.CXXFlags...)
				compileCmds = append(compileCmds, CompileCommand{Directory: srcDir, Arguments: append([]string{os.Getenv("CXX")},
					compileArgs...), File: path})
				fallthrough
			case ".hpp":
				isCXX = true
				jsonArgs := append(slices.Clip(options.CommandOnlyFlags.CXXFlags), options.BuildFlags.CXXFlags...)
				jsonArgs = append(jsonArgs, "-c", path)
				jsonCmds = append(jsonCmds,
					CompileCommand{Directory: srcDir, Arguments: append([]string{os.Getenv("CXX")},
						jsonArgs...), File: path})
			}
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	ok := atomic.Bool{}
	ok.Store(true)
	wg := sync.WaitGroup{}
	for i, c := range compileCmds {
		obj := filepath.Join(buildDir, objs[i])
		if err := os.MkdirAll(filepath.Dir(obj), 0o755); err != nil {
			return nil, err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if out, err := toolchain.RunCombinedOutput(c.Arguments[0], append(c.Arguments[1:], "-o", obj, "-c", filepath.Join(c.Directory, c.File))...); err != nil {
				debug.EPrintf("%s", out)
				ok.Store(false)
			}
		}()
	}
	wg.Wait()
	if !ok.Load() {
		return nil, debug.Errorf("Failed to build %q", srcDir)
	}
	if err := os.MkdirAll(filepath.Dir(outFile), 0o755); err != nil {
		return nil, err
	}
	switch options.Type {
	case BuildTypeStaticLibrary:
		args := []string{"rcs", outFile}
		if err := toolchain.RunDir(buildDir, os.Getenv("AR"), append(args, objs...)...); err != nil {
			return nil, err
		}
	case BuildTypeSharedLibrary:
		args := []string{"-o", outFile, "-rdynamic", "-shared"}
		args = append(args, objs...)
		args = append(args, options.LDFlags...)
		if isCXX {
			if err := toolchain.RunDir(buildDir, os.Getenv("CXX"), args...); err != nil {
				return nil, err
			}
		} else {
			if err := toolchain.RunDir(buildDir, os.Getenv("CC"), args...); err != nil {
				return nil, err
			}
		}
	case BuildTypeExecuteable:
		args := []string{"-o", outFile}
		args = append(args, objs...)
		args = append(args, options.LDFlags...)
		if isCXX {
			if err := toolchain.RunDir(buildDir, os.Getenv("CXX"), args...); err != nil {
				return nil, err
			}
		} else {
			if err := toolchain.RunDir(buildDir, os.Getenv("CC"), args...); err != nil {
				return nil, err
			}
		}
	}
	return jsonCmds, nil
}
