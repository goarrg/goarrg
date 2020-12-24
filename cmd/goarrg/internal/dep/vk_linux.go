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
	"runtime"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/debug"
)

func init() {
	switch base.GOOS() {
	case runtime.GOOS:
		deps["vk"] = dep{
			preBuild: func() {
				vkPath := os.Getenv("VULKAN_SDK")

				if vkPath == "" {
					if s, err := os.Stat("/usr/include/vulkan"); err == nil && s.IsDir() {
						return
					}

					os.Remove(usrData + "/lib/pkgconfig/vulkan.pc")
					flagDisableVK = true
					debug.LogI("Vulkan headers not found and VULKAN_SDK unset, disabling Vulkan")
					return
				}

				if f, err := os.Create(usrData + "/lib/pkgconfig/vulkan.pc"); err != nil {
					panic(err)
				} else {
					defer f.Close()

					_, err = f.Write([]byte("Name: vulkan" +
						"\nDescription:" +
						"\nVersion:" +
						"\nRequires:" +
						"\nConflicts:" +
						"\nLibs: -L" + vkPath + "/lib" +
						"\nCflags: -I" + vkPath + "/include" +
						"\n"))

					if err != nil {
						panic(err)
					}
				}
			},
		}
	default:
		deps["vk"] = dep{
			install: func() {
				if err := os.MkdirAll(usrData+"/include", 0o755); err != nil {
					panic(err)
				}

				if s, err := os.Stat("/usr/include/vulkan"); err == nil && s.IsDir() {
					err := os.Symlink("/usr/include/vulkan", usrData+"/include/vulkan")

					if err != nil {
						panic(err)
					}
				}
			},
			preBuild: func() {
				vkPath := os.Getenv("VULKAN_SDK")

				if vkPath == "" {
					if s, err := os.Stat(usrData + "/include/vulkan"); err != nil || !s.IsDir() {
						os.Remove(usrData + "/lib/pkgconfig/vulkan.pc")
						flagDisableVK = true
						debug.LogI("Vulkan headers not found and VULKAN_SDK unset, disabling Vulkan")
						return
					}

					vkPath = usrData
				}

				if f, err := os.Create(usrData + "/lib/pkgconfig/vulkan.pc"); err != nil {
					panic(err)
				} else {
					defer f.Close()

					_, err = f.Write([]byte("Name: vulkan" +
						"\nDescription:" +
						"\nVersion:" +
						"\nRequires:" +
						"\nConflicts:" +
						"\nLibs: -L" + vkPath + "/lib" +
						"\nCflags: -I" + vkPath + "/include" +
						"\n"))

					if err != nil {
						panic(err)
					}
				}
			},
		}
	}
}
