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
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/debug"
)

func init() {
	switch base.GOOS() {
	case "windows":
		deps["vk"] = dep{
			preBuild: func() {
				vkPath := os.Getenv("VULKAN_SDK")

				if vkPath == "" {
					os.Remove(filepath.Join(usrData, "lib", "pkgconfig", "vulkan.pc"))
					flagDisableVK = true
					debug.LogI("VULKAN_SDK unset, disabling vk")
					return
				}

				vkPath = filepath.ToSlash(vkPath)

				if f, err := os.Create(filepath.Join(usrData, "lib", "pkgconfig", "vulkan.pc")); err != nil {
					panic(err)
				} else {
					defer f.Close()

					_, err = f.Write([]byte("prefix=" + filepath.ToSlash(usrData) +
						"\nexec_prefix=${prefix}" +
						"\nlibdir=${exec_prefix}/lib" +
						"\nincludedir=${prefix}/include" +
						"\n" +
						"\nName: vulkan" +
						"\nDescription:" +
						"\nVersion:" +
						"\nRequires:" +
						"\nConflicts:" +
						"\nLibs: -L${libdir} -L" + vkPath + "/Lib" +
						"\nCflags: -I${includedir} -I" + vkPath + "/Include" +
						"\n"))

					if err != nil {
						panic(err)
					}
				}
			},
		}
	}
}
