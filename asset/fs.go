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

package asset

import (
	"io/fs"
	"path/filepath"
)

type FileSystem struct {
	dir string
}

var FS = &FileSystem{dir: "./"}

func (f *FileSystem) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}
	return Load(filepath.Join(f.dir, name))
}

func DirFS(dir string) *FileSystem {
	return &FileSystem{dir: dir}
}
