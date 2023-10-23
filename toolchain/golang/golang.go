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

package golang

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"goarrg.com/debug"
	"goarrg.com/toolchain"

	"golang.org/x/tools/go/packages"
)

type Config struct {
	Target toolchain.Target
}

var goEnvOnce = sync.Once{}

func Setup(c Config) {
	debug.IPrintf("Setting up golang")

	{
		gocache := filepath.Join(toolchain.WorkingModuleDataDir(), "gocache")
		if err := os.MkdirAll(gocache, 0o755); err != nil {
			panic(err)
		}
		toolchain.EnvRegister("GOCACHE", gocache)
	}

	if !ValidTarget(c.Target) {
		panic(debug.Errorf("Unknown os/arch combo: %s", c.Target))
	}

	toolchain.EnvSet("GOOS", c.Target.OS)
	toolchain.EnvSet("GOARCH", c.Target.Arch)

	if !CgoSupported(c.Target) {
		debug.WPrintf("cgo unsupported on target: %s", c.Target)
	} else {
		toolchain.EnvSet("CGO_ENABLED", "1")
	}

	if c.Target.CrossCompiling() {
		debug.IPrintf("Detected cross compiling target: %s", c.Target)
	}

	goEnvOnce.Do(func() {
		goEnv, err := toolchain.RunOutput("go", "env", "-json")
		if err != nil {
			panic(err)
		}
		goEnvMap := map[string]string{}
		if err := json.Unmarshal(goEnv, &goEnvMap); err != nil {
			panic(err)
		}
		// merge go's env with ours but ignore go's value if we've set it
		for k, v := range goEnvMap {
			toolchain.EnvRegister(k, v)
		}
	})
}

func AppendTag(args []string, tag string) []string {
	haveTagsArg := false
	for i, arg := range args {
		if strings.HasPrefix(arg, "-tags") {
			arg = strings.TrimPrefix(arg, "-tags")
			if strings.HasPrefix(arg, "=") {
				arg = strings.TrimPrefix(arg, "=")
				args[i] = "-tags=" + strings.ReplaceAll(arg, " ", ",") + "," + tag
			} else {
				arg = args[i+1]
				args[i+1] = strings.ReplaceAll(arg, " ", ",") + "," + tag
			}
			haveTagsArg = true
			break
		}
	}

	if !haveTagsArg {
		args = append(args, "-tags="+tag)
	}

	return args
}

/*
CallersPackage is a convenience function that returns the packages.Package of the
package of the function that called CallersPackage.
*/
func CallersPackage(mode packages.LoadMode) *packages.Package {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(debug.Errorf("Failed runtime.Caller"))
	}

	p, err := packages.Load(&packages.Config{Mode: mode}, filepath.Dir(file))
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to load package of caller"))
	}
	if len(p) == 0 {
		panic(debug.Errorf("No go package for caller?!?"))
	}

	// there should only be one
	return p[0]
}

/*
CallersModule is a convenience function that returns the packages.Module of the
package of the function that called CallersModule. It is effectively
CallersPackage(packages.NeedModule).Module.
*/
func CallersModule() *packages.Module {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(debug.Errorf("Failed runtime.Caller"))
	}

	p, err := packages.Load(&packages.Config{Mode: packages.NeedModule}, filepath.Dir(file))
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to load package of caller"))
	}
	if len(p) == 0 {
		panic(debug.Errorf("No go package for caller?!?"))
	}

	// if len is > 1 the module should still all be the same
	return p[0].Module
}

var shouldCleanCache bool

/*
SetShouldCleanCache signals that the go cache should be cleared. This is useful
for packages to signal that the go cache should be cleared due to how things
work with cgo.
*/
func SetShouldCleanCache() { shouldCleanCache = true }

/*
ShouldCleanCache returns true when SetShouldCleanCache() has been called.
*/
func ShouldCleanCache() bool { return shouldCleanCache }

/*
CleanCache is a convenience function that cleans the go cache.
*/
func CleanCache() {
	if err := toolchain.Run("go", "clean", "-cache", "-testcache"); err != nil {
		panic(err)
	}
	// prevent redundant cleans
	shouldCleanCache = false
}
