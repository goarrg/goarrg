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

package gcc

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
)

func gccArch(goarch string) string {
	switch goarch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i686"
	}
	return ""
}

func gccABI(goos string) string {
	switch goos {
	case "linux":
		return "linux-gnu"
	case "windows":
		return "w64-mingw32"
	}
	return ""
}

func gccTarget(target toolchain.Target) string {
	arch := gccArch(target.Arch)
	abi := gccABI(target.OS)
	if arch == "" || abi == "" {
		return ""
	}
	return arch + "-" + abi
}

func gccEnv(target toolchain.Target) map[string]string {
	gccTarget := gccTarget(target) + "-"
	if !target.CrossCompiling() {
		gccTarget = "" // only use target specific filenames when crosscompiling
	}
	m := map[string]string{
		"CC":  gccTarget + "gcc",
		"CXX": gccTarget + "g++",
		"AR":  gccTarget + "gcc-ar",
		"RC":  "",
	}

	switch target.OS {
	case "windows":
		m["RC"] = gccTarget + "windres"
	}
	return m
}

type Config struct {
	Target toolchain.Target
}

func Setup(c Config) {
	debug.IPrintf("Setting up gcc toolchain")

	if gccTarget := gccTarget(c.Target); gccTarget == "" {
		panic(debug.Errorf("GCC target specific toolchain filenames not known for: %s", c.Target))
	} else if !Installed(c.Target) {
		panic(debug.Errorf("GCC not found for target: %s", gccTarget))
	}

	gccEnv := gccEnv(c.Target)
	for k, v := range gccEnv {
		toolchain.EnvSet(k, v)
	}
}

func FindHeader(target toolchain.Target, header string) (string, error) {
	if gccTarget := gccTarget(target); gccTarget == "" {
		return "", debug.Errorf("GCC target specific toolchain filenames not known for: %s", target)
	}

	gccEnv := gccEnv(target)
	ex := exec.Command(gccEnv["CC"], "-M", "-E", "-")
	ex.Env = os.Environ()
	ex.Stdin = bytes.NewReader([]byte("#include<" + header + ">"))
	out, err := ex.CombinedOutput()
	if err != nil {
		return "", debug.Errorf("Failed to find %q using %q: %q", header, gccEnv["CC"], string(out))
	}

	for _, s := range strings.Fields(string(out)) {
		if strings.Contains(s, header) {
			return s, nil
		}
	}

	// should never be here
	panic(debug.Errorf("%q found %q but unable to find header in output: %q", gccEnv["CC"], header, string(out)))
}

func Installed(target toolchain.Target) bool {
	gccEnv := gccEnv(target)
	for _, v := range gccEnv {
		if v == "" {
			continue
		}
		if _, err := exec.LookPath(v); err != nil {
			return false
		}
	}
	return true
}
