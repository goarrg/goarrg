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
	"errors"
	"flag"
	"runtime"
	"strings"
)

var flagTarget = runtime.GOOS + "_" + runtime.GOARCH

const flagTargetUsage = `Build target, must be either ${GOOS} or ${GOOS}_${GOARCH}.
If ${GOARCH} is unspecified, will default to runtime.GOARCH.`

var (
	targetOS   = runtime.GOOS
	targetArch = runtime.GOARCH
)

func AddFlags(f *flag.FlagSet) {
	f.Func("target", flagTargetUsage, func(s string) error {
		if s == "" {
			return errors.New("Target is empty")
		}

		if osarch := strings.SplitN(s, "_", 2); len(osarch) < 2 {
			targetOS = osarch[0]
			targetArch = runtime.GOARCH
		} else {
			targetOS = osarch[0]
			targetArch = osarch[1]
		}

		flagTarget = targetOS + "_" + targetArch
		if ValidPlatform(flagTarget) {
			return nil
		}

		return errors.New("Unknown target")
	})
}

func Target() string {
	return flagTarget
}

func TargetOS() string {
	return targetOS
}

func TargetArch() string {
	return targetArch
}
