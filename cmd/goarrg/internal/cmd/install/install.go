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

package install

import (
	"goarrg.com/cmd/goarrg/internal/cgodep"
	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
)

const short = `Install C dependencies and goarrg-config files for use with "#cgo pkg-config: [cgodep]".
Installed dependencies will be automatically redownloaded/rebuilt as needed.`

const long = short + `

NOTE:
- Only available on targets supporting cgo.
- goarrg-config files are not readable by pkg-config.
`

var CMD = &cmd.CMD{
	Run:   run,
	Name:  "install",
	Short: short,
	Long:  long,
	CMDs:  map[string]*cmd.CMD{},
}

func init() {
	for _, d := range cgodep.List() {
		d := d
		CMD.CMDs[d.Name] = &cmd.CMD{
			Name: d.Name,
			Run: func(args []string) bool {
				if len(args) != 0 {
					debug.LogE("Invalid args: %q", args)
					return false
				}

				toolchain.Setup()
				d.Install()

				// we have to clean the cache everytime we change a external C dependency
				// cause the cache would not detect those changes and so would use the
				// old cached object files.
				if err := exec.Run("go", "clean", "-cache"); err != nil {
					panic(err)
				}
				return true
			},
			Short: d.Short,
			Long:  d.Long,
		}

		if d.TargetSpecific {
			toolchain.AddFlags(&CMD.CMDs[d.Name].Flags)
		}
	}
}

func run(args []string) bool {
	debug.LogE("Please specify what do you want to install")
	if len(args) != 0 {
		debug.LogE("Invalid args: %q", args)
	}
	return false
}
