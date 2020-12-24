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

package clean

import (
	"os"
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/cmd/goarrg/internal/exec"
)

var CMD = &base.CMD{
	Run: func(args []string) bool {
		if len(args) > 0 {
			return false
		}

		{
			gocache, err := filepath.Abs(filepath.Join(".", ".goarrg", "gocache"))

			if err != nil {
				panic(err)
			}

			if err := os.MkdirAll(gocache, 0o755); err != nil {
				panic(err)
			}

			if err := os.Setenv("GOCACHE", gocache); err != nil {
				panic(err)
			}
		}

		if err := exec.Run("go", "clean", "-cache", "-testcache", "./..."); err != nil {
			panic(err)
		}

		return true
	},
	Name:  "clean",
	Short: "Clears build cache and current directory of genereated files",
	Long:  "",
	CMDs: map[string]*base.CMD{
		"yourself": {
			Run: func(args []string) bool {
				if len(args) > 0 {
					return false
				}

				if err := exec.Run("go", "clean", "-cache", "-testcache", "goarrg.com/..."); err != nil {
					panic(err)
				}

				return true
			},
			Name:  "yourself",
			Short: "Clears build cache and goarrg sources of genereated files",
			Long:  "",
		},
	},
}
