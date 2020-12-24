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

package cgo

import (
	"os"
	"runtime"

	"goarrg.com/cmd/goarrg/internal/base"
)

func init() {
	err := os.Setenv("CGO_ENABLED", "1")

	if err != nil {
		panic(err)
	}

	if runtime.GOOS == base.GOOS() {
		return
	}

	switch base.GOOS() {
	case "windows":
		err := os.Setenv("CC", GCCArch()+"-w64-mingw32-gcc")

		if err != nil {
			panic(err)
		}

		err = os.Setenv("CXX", GCCArch()+"-w64-mingw32-g++")

		if err != nil {
			panic(err)
		}
	default:
		panic("No support for target os: " + base.GOOS())
	}
}

func GCCArch() string {
	switch base.GOARCH() {
	case "amd64":
		return "x86_64"

	case "386":
		return "i686"

	default:
		panic("No support for target arch: " + base.GOARCH())
	}
}
