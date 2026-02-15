/*
Copyright 2022 The goARRG Authors.

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

package goarrg

import (
	"fmt"
	"path/filepath"

	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cgodep"
	"goarrg.com/toolchain/git"
	"goarrg.com/toolchain/golang"
	"goarrg.com/toolchain/web"
)

const (
	vkHeadersURL = "https://github.com/KhronosGroup/Vulkan-Headers.git"
	vkDocsURL    = "https://github.com/KhronosGroup/Vulkan-Docs.git"
	vkBuild      = "-goarrg0"
)

type vkRepo struct {
	URL    string
	Tag    string
	Commit string
}

func getVkRepoData(os string, version string) map[string]vkRepo {
	if version == "" {
		version = "latest"
	}
	switch os {
	case "linux", "windows":
	case "darwin":
		os = "mac"
	default:
		panic(fmt.Sprintf("Unsupported os: %s", os))
	}
	type config struct {
		Repos map[string]vkRepo
	}
	cfg := config{}
	err := web.GetJSON(fmt.Sprintf("https://vulkan.lunarg.com/sdk/config/%s/%s/config.json", version, os), &cfg)
	if err != nil {
		panic(err)
	}
	return cfg.Repos
}

func installVkHeaders(tag string) error {
	installDir := cgodep.InstallDir("vulkan-headers", toolchain.Target{}, toolchain.BuildRelease)
	var ref git.Ref
	{
		refs, err := git.GetRemoteTags(vkHeadersURL, tag)
		if err != nil {
			return err
		}
		ref = refs[0]
	}
	buildID := ref.Name + "-" + ref.Hash + vkBuild
	if cgodep.ReadVersion(installDir) == buildID {
		return nil
	}
	err := git.CloneOrFetch(vkHeadersURL, installDir, ref)
	if err != nil {
		return err
	}
	golang.SetShouldCleanCache()
	return cgodep.WriteMetaFile("vulkan-headers", toolchain.Target{}, toolchain.BuildRelease, cgodep.Meta{
		Version: buildID, Flags: cgodep.Flags{
			CFlags: []string{"-I" + filepath.Join(installDir, "include")},
		},
	})
}

func installVkDocs(tag string) error {
	installDir := cgodep.InstallDir("vulkan-docs", toolchain.Target{}, toolchain.BuildRelease)
	var ref git.Ref
	{
		refs, err := git.GetRemoteTags(vkDocsURL, tag)
		if err != nil {
			return err
		}
		ref = refs[0]
	}
	buildID := ref.Name + "-" + ref.Hash + vkBuild
	if cgodep.ReadVersion(installDir) == buildID {
		return nil
	}
	err := git.CloneOrFetch(vkDocsURL, installDir, ref)
	if err != nil {
		return err
	}
	return cgodep.WriteMetaFile("vulkan-docs", toolchain.Target{}, toolchain.BuildRelease, cgodep.Meta{
		Version: buildID,
	})
}
