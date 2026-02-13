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
	"runtime"

	"goarrg.com/build"
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
