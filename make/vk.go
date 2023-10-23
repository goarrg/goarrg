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
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cgodep"
	"goarrg.com/toolchain/golang"
)

const (
	vkHeadersBuild = "-goarrg0"
)

type VkHeadersConfig struct {
	Install bool

	Branch string
}

func installVkHeaders(c VkHeadersConfig) error {
	installDir := cgodep.InstallDir("vulkan-headers", toolchain.Target{}, toolchain.BuildRelease)
	var headHash string
	if c.Branch == "" {
		branches, err := toolchain.RunOutput("git", "ls-remote", "--heads", "--sort=-version:refname", "https://github.com/KhronosGroup/Vulkan-Headers.git", "sdk-*")
		if err != nil {
			return err
		}
		headHash = string(branches[:bytes.IndexAny(branches, " \t")])
		i := bytes.Index(branches, []byte("sdk-"))
		j := bytes.Index(branches[i:], []byte("\n"))
		c.Branch = strings.TrimSpace(string(branches[i : i+j]))
	} else {
		output, err := toolchain.RunOutput("git", "ls-remote", "--heads", "--sort=-version:refname", "https://github.com/KhronosGroup/Vulkan-Headers.git", c.Branch)
		if err != nil {
			return err
		}
		headHash = string(output[:bytes.IndexAny(output, " \t")])
	}

	{
		version := cgodep.ReadVersion(installDir)
		// search backwards cause branch names have a "-" in them
		i := strings.LastIndex(version, "-")
		if i > 0 {
			j := strings.LastIndex(version[:i], "-")
			if j > 0 {
				branch := version[:j]
				hash := version[j+1 : i]
				build := version[i:]

				if branch == c.Branch && hash == headHash && build == vkHeadersBuild {
					return nil
				}
			}
		}
	}

	if stat, err := os.Stat(filepath.Join(installDir, ".git")); err != nil || !stat.IsDir() {
		os.RemoveAll(installDir)
		if err := os.MkdirAll(installDir, 0o755); err != nil {
			return err
		}
		if err := toolchain.Run("git", "clone", "--branch", c.Branch, "--depth=1", "https://github.com/KhronosGroup/Vulkan-Headers.git", installDir); err != nil {
			return err
		}
	}
	if err := toolchain.RunDir(installDir, "git", "fetch", "origin", "--depth=1", "refs/heads/"+c.Branch); err != nil {
		return err
	}
	if err := toolchain.RunDir(installDir, "git", "-c", "advice.detachedHead=false", "checkout", "FETCH_HEAD"); err != nil {
		return err
	}
	if err := toolchain.RunDir(installDir, "git", "clean", "-q", "-f", "-d"); err != nil {
		return err
	}
	hash, err := toolchain.RunDirOutput(installDir, "git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	golang.SetShouldCleanCache()
	return cgodep.WriteMetaFile("vulkan-headers", toolchain.Target{}, toolchain.BuildRelease, cgodep.Meta{
		Version: c.Branch + "-" + strings.TrimSpace(string(hash)) + vkHeadersBuild, Flags: cgodep.Flags{
			CFlags: []string{"-I" + filepath.Join(installDir, "include")},
		},
	})
}
