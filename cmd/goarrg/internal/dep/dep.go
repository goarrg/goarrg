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

package dep

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goarrg.com/cmd/goarrg/internal/archive"
	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/remote"
	"goarrg.com/debug"
)

type dep struct {
	file, url, version string
	verify             func(io.ReadSeeker) error
	install            func()
	preBuild           func()
}

type depMeta struct {
	Version map[string]string
}

var deps = map[string]dep{}

var usrData = filepath.Join(base.UsrData(), "deps", base.GOOS(), base.GOARCH())
var usrCache = filepath.Join(base.UsrCache(), "deps")

func init() {
	pkg := filepath.Join(usrData, "lib", "pkgconfig")

	if oldPkg := os.Getenv("PKG_CONFIG_PATH"); oldPkg == "" {
		os.Setenv("PKG_CONFIG_PATH", pkg)
	} else {
		os.Setenv("PKG_CONFIG_PATH", pkg+string(filepath.ListSeparator)+oldPkg)
	}
}

func depCheck() {
	f, err := os.Open(filepath.Join(usrData, "dep.json"))

	if err != nil {
		flagDep = true
		return
	}

	defer f.Close()
	data, err := ioutil.ReadAll(f)

	if err != nil {
		panic(err)
	}

	meta := depMeta{}
	if err := json.Unmarshal(data, &meta); err != nil {
		panic(err)
	}

	for k := range deps {
		if meta.Version[k] == deps[k].version {
			continue
		}

		flagDep = true
		return
	}
}

func depDownload(name string, dep dep) {
	retried := false
	srcDir := filepath.Join(usrCache, name)

	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		panic(err)
	}

	if files, err := ioutil.ReadDir(srcDir); err != nil {
		panic(err)
	} else {
		for _, file := range files {
			if file.Name() != dep.file {
				os.RemoveAll(filepath.Join(srcDir, file.Name()))
			}
		}
	}

	src, err := os.Open(filepath.Join(srcDir, dep.file))

	if err == nil {
		if dep.verify == nil || dep.verify(src) == nil {
			defer src.Close()
			goto extract
		} else {
			src.Close()
		}
	}

	src, err = os.Create(filepath.Join(srcDir, dep.file))

	if err != nil {
		panic(err)
	}

	defer src.Close()

retry:

	debug.LogI("Getting dependency %s", name)

	if err := remote.Get(dep.url, src); err != nil {
		panic(err)
	}

	if dep.verify != nil {
		debug.LogI("Verifying dependency %s", name)
		if err = dep.verify(src); err != nil {
			if !retried {
				debug.LogE("Failed to verify %s, redownloading", dep.file)
				if err := src.Truncate(0); err != nil {
					panic(err)
				}
				retried = true
				goto retry
			}
			panic(err)
		}
	}

extract:

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	if strings.HasSuffix(dep.file, ".tar.gz") {
		if err := archive.ExtractHere(src); err != nil {
			if !retried {
				debug.LogE("Failed to extract %s, redownloading", dep.file)
				if err := src.Truncate(0); err != nil {
					panic(err)
				}
				retried = true
				goto retry
			}
			panic(err)
		}
	}
}

func depPrebuild() {
	for name, dep := range deps {
		if dep.preBuild != nil {
			debug.LogI("Doing pre-build step for dependency %s", name)
			dep.preBuild()
			debug.LogI("Done pre-build step for dependency %s", name)
		}
	}
}

func Build() {
	if depCheck(); !flagDep {
		if !flagNoDep {
			depPrebuild()
		}
		return
	}

	if flagNoDep {
		if err := os.Chtimes(usrData, time.Now(), time.Now()); err != nil {
			panic(err)
		}
		return
	}

	//clear build cache to prevent linking outdated C libs/headers
	if err := exec.Run("go", "clean", "-cache", "-testcache"); err != nil {
		panic(err)
	}

	if err := os.RemoveAll(usrData); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(usrData, "lib", "pkgconfig"), 0o755); err != nil {
		panic(err)
	}

	meta := depMeta{
		Version: make(map[string]string),
	}

	for name, dep := range deps {
		buildDir, err := ioutil.TempDir("", name)

		if err != nil {
			panic(err)
		}

		if err := os.Chdir(buildDir); err != nil {
			panic(err)
		}

		if dep.file != "" {
			depDownload(name, dep)
		}

		if dep.install != nil {
			debug.LogI("Installing dependency %s", name)
			dep.install()
			meta.Version[name] = dep.version
			debug.LogI("Done installing dependency %s", name)
		}

		os.RemoveAll(buildDir)
	}

	data, err := json.Marshal(meta)

	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(usrData, "dep.json"), data, 0o644); err != nil {
		panic(err)
	}

	base.ResetCWD()
	depPrebuild()
}
