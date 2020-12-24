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
	"os"

	"goarrg.com/debug"
)

func init() {
	deps["gl"] = dep{
		preBuild: func() {
			if s, err := os.Stat("/usr/include/GL"); err == nil && s.IsDir() {
				if s, err := os.Stat("/usr/include/GL/glu.h"); err == nil && !s.IsDir() {
					return
				}
			}

			flagDisableGL = true
			debug.LogI("OpenGL headers not found, disabling OpenGL")
		},
	}
}
