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

package cc

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
)

type Compiler int

const (
	CompilerGCC = iota
	CompilerClang
)

func (c Compiler) String() string {
	switch c {
	case CompilerGCC:
		return "gcc"
	case CompilerClang:
		return "clang"
	}

	panic(debug.Errorf("Unknown compiler value: %d", c))
}

func gnuArch(goarch string) string {
	switch goarch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i686"
	}
	return ""
}

func gnuABI(goos string) string {
	switch goos {
	case "linux":
		return "linux-gnu"
	case "windows":
		return "w64-mingw32"
	}
	return ""
}

func gnuPrefix(target toolchain.Target) string {
	arch := gnuArch(target.Arch)
	abi := gnuABI(target.OS)
	if arch == "" || abi == "" {
		panic(debug.Errorf("No known compiler prefix for target: %q", target))
	}
	return arch + "-" + abi + "-"
}

func gccEnv(target toolchain.Target) map[string]string {
	prefix := ""
	if target.CrossCompiling() {
		prefix = gnuPrefix(target)
	}
	m := map[string]string{
		"CC":  prefix + "gcc",
		"CXX": prefix + "g++",
		"AR":  prefix + "gcc-ar",
		"RC":  "",
	}
	switch target.OS {
	case "windows":
		m["RC"] = prefix + "windres"
	}
	return m
}

func clangEnv(target toolchain.Target) map[string]string {
	prefix := ""
	if target.CrossCompiling() {
		prefix = gnuPrefix(target)
	}
	m := map[string]string{
		"CC":  prefix + "clang",
		"CXX": prefix + "clang++",
		"AR":  prefix + "llvm-ar",
		"RC":  "",
	}
	switch target.OS {
	case "windows":
		m["RC"] = prefix + "windres"
	}
	return m
}

func compilerEnv(c Config) map[string]string {
	switch c.Compiler {
	case CompilerGCC:
		return gccEnv(c.Target)
	case CompilerClang:
		return clangEnv(c.Target)
	}
	panic(debug.Errorf("Unknown compiler value: %d", c.Compiler))
}

type Config struct {
	Compiler Compiler
	Target   toolchain.Target
}

func Setup(c Config) {
	debug.IPrintf("Setting up C/C++ toolchain")

	if !Installed(c) {
		panic(debug.Errorf("Compiler not installed for Config: %+v", c))
	}

	env := compilerEnv(c)
	for k, v := range env {
		toolchain.EnvSet(k, v)
	}
}

func FindHeader(cfg Config, header string) (string, error) {
	cc := toolchain.EnvGet("CC")
	if cfg != (Config{}) {
		env := compilerEnv(cfg)
		cc = env["CC"]
	}

	ex := exec.Command(cc, "-M", "-E", "-")
	ex.Env = os.Environ()
	ex.Stdin = bytes.NewReader([]byte("#include<" + header + ">"))
	out, err := ex.CombinedOutput()
	if err != nil {
		return "", debug.Errorf("Failed to find %q using %q: %q", header, cc, string(out))
	}

	for _, s := range strings.Fields(string(out)) {
		if strings.Contains(s, header) {
			return s, nil
		}
	}

	// should never be here
	panic(debug.Errorf("%q found %q but unable to find header in output: %q", cc, header, string(out)))
}

func Installed(c Config) bool {
	found := true
	for _, v := range compilerEnv(c) {
		if v == "" {
			continue
		}
		if _, err := exec.LookPath(v); err != nil {
			debug.VPrintf("Failed to find: %q", v)
			found = false
		}
	}
	return found
}
