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
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
)

/*
Install installs the helper program ./cmd/cgodep-config as a replacement for
PKG_CONFIG. Currently it is only able to use dependencies installed through
this package. It might be possible in the future for it to read PKG_CONFIG files.
*/
func Install() {
	tool := "cgodep-config"
	toolFile := filepath.Join(toolchain.ToolsDir(), tool)

	if runtime.GOOS == "windows" {
		toolFile += ".exe"
	}
	if err := toolchain.RunEnv(append(os.Environ(), "GOBIN="+toolchain.ToolsDir(), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH), "go", "install",
		"goarrg.com/cmd/"+tool); err != nil {
		panic(err)
	}
	toolchain.EnvRegister("CGODEP_PATH", filepath.Join(toolchain.WorkingModuleDataDir(), "cgodep"))
	toolchain.EnvSet("PKG_CONFIG", toolFile)
}

/*
Flags used by cgodep-config for when go invokes ${PKG_CONFIG}
*/
type Flags struct {
	CFlags  []string
	LDFlags []string

	/*
		StaticLDFlags are the LDFlags that would be passed to cgo if pkg-config
		was executed with the "--static" flag. However, unlike pkg-config,
		goarrg would not combine LDFlags and StaticLDFlags.
	*/
	StaticLDFlags []string
}

type Meta struct {
	Version string
	Flags   Flags
}

const metaFileName = "goarrg_cgodep.json"

/*
DataDir returns the path where all cgodep related data should be stored.
*/
func DataDir() string {
	searchList := filepath.SplitList(os.Getenv("CGODEP_PATH"))

	// select first non empty path, else use default
	for i := range searchList {
		if searchList[i] != "" {
			return searchList[i]
		}
	}

	return filepath.Join(toolchain.WorkingModuleDataDir(), "cgodep")
}

/*
InstallDir returns the path where the data for "name" should be stored. It is
specific to target and build. So there is no concern of one target/build overriding another.
As a special case, if target is the zero value, it is only given a path specific to name.
This allows for header only dependencies without having to have multiple copies for every
target/build.
*/
func InstallDir(name string, target toolchain.Target, build toolchain.Build) string {
	if (target == toolchain.Target{}) {
		return filepath.Join(DataDir(), name)
	}
	return filepath.Join(DataDir(), name, target.String(), build.String())
}

/*
SetActiveBuild is used to select which build should be used for building. The
active build is independent for each target.
*/
func SetActiveBuild(name string, target toolchain.Target, build toolchain.Build) error {
	if (target == toolchain.Target{}) {
		return nil
	}

	installedMetaFile := filepath.Join(InstallDir(name, target, build), metaFileName)
	fIn, err := os.Open(installedMetaFile)
	if err != nil {
		return debug.ErrorWrapf(err, "Failed to read: %q", installedMetaFile)
	}
	defer fIn.Close()

	resolveMetaFile := filepath.Join(filepath.Join(DataDir(), name, target.String()), metaFileName)
	fOut, err := os.OpenFile(resolveMetaFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return debug.ErrorWrapf(err, "Failed to create: %q", resolveMetaFile)
	}
	defer fOut.Close()

	_, err = io.Copy(fOut, fIn)
	return debug.ErrorWrapf(err, "Failed to write metafile")
}

/*
WriteMetaFile writes the metadata for name to InstallDir.
*/
func WriteMetaFile(name string, target toolchain.Target, build toolchain.Build, m Meta) error {
	j, err := json.Marshal(m)
	if err != nil {
		return debug.ErrorWrapf(err, "Failed to marshal metafile")
	}
	installDir := InstallDir(name, target, build)
	metaFile := filepath.Join(installDir, metaFileName)
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(metaFile, j, 0o644); err != nil {
		return debug.ErrorWrapf(err, "Failed to write metafile")
	}
	return SetActiveBuild(name, target, build)
}

/*
ReadMetaFile returns the metadata located at dir and nil or nil and error on error.
*/
func ReadMetaFile(dir string) (Meta, error) {
	metaFile := filepath.Join(dir, metaFileName)
	j, err := os.ReadFile(metaFile)
	if err != nil {
		return Meta{}, debug.ErrorWrapf(err, "Failed to read: %q", metaFile)
	}
	m := Meta{}
	err = json.Unmarshal(j, &m)
	return m, debug.ErrorWrapf(err, "Failed to unmarshal metafile")
}

/*
ReadVersion is a convenience function that is basically calling ReadMetaFile
and returning the version or panics on error unless error is os.ErrNotExist.
*/
func ReadVersion(dir string) string {
	m, err := ReadMetaFile(dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return m.Version
}
