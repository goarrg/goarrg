/*
Copyright 2023 The goARRG Authors.

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

package git

import (
	"os"
	"path/filepath"
	"strings"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
)

type Ref struct {
	Hash string
	Name string
}

func GetRemoteTags(url, pattern string) ([]Ref, error) {
	var out []Ref
	ls, err := toolchain.RunOutput("git", "ls-remote", "--refs", "--tags", "--sort=-version:refname", url, pattern)
	if err != nil {
		return nil, err
	}
	lsStr := strings.TrimSpace(string(ls))
	if lsStr == "" {
		return nil, debug.Errorf("GetRemoteTags: No results")
	}
	refs := strings.Split(lsStr, "\n")
	for _, r := range refs {
		hash := r[:strings.IndexAny(r, " \t")]
		ref := strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(r, hash)), "refs/tags/")
		out = append(out, Ref{Hash: hash, Name: ref})
	}
	return out, nil
}

func GetRemoteHeads(url, pattern string) ([]Ref, error) {
	var out []Ref
	ls, err := toolchain.RunOutput("git", "ls-remote", "--heads", "--sort=-version:refname", url, pattern)
	if err != nil {
		return nil, err
	}
	lsStr := strings.TrimSpace(string(ls))
	if lsStr == "" {
		return nil, debug.Errorf("GetRemoteHeads: No results")
	}
	refs := strings.Split(lsStr, "\n")
	for _, r := range refs {
		hash := r[:strings.IndexAny(r, " \t")]
		ref := strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(r, hash)), "refs/heads/")
		out = append(out, Ref{Hash: hash, Name: ref})
	}
	return out, nil
}

func CloneOrFetch(url, dir string, ref Ref) error {
	if stat, err := os.Stat(filepath.Join(dir, ".git")); err != nil || !stat.IsDir() {
		os.RemoveAll(dir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := toolchain.Run("git", "-c", "advice.detachedHead=false",
			"clone", "--branch", ref.Name, "--depth=1", url, dir); err != nil {
			return err
		}
	}
	if err := toolchain.RunDir(dir, "git", "fetch", "origin", "--depth=1", ref.Name); err != nil {
		return err
	}
	if err := toolchain.RunDir(dir, "git", "-c", "advice.detachedHead=false", "checkout", ref.Hash); err != nil {
		return err
	}
	if err := toolchain.RunDir(dir, "git", "clean", "-q", "-f", "-d"); err != nil {
		return err
	}
	return nil
}
