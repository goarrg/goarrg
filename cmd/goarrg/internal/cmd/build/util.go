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

package build

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/debug"
)

//nolint:deadcode,unused
func generatedFileName(name, ext string) string {
	return name + "_generated_" + base.GOOS() + "_" + base.GOARCH() + ext
}

//nolint:unused
func copyFile(src, dest string) error {
	from, err := os.Open(src)

	if err != nil {
		return debug.ErrorWrap(err, "Open file failed")
	}

	defer from.Close()
	to, err := os.Create(dest)

	if err != nil {
		return debug.ErrorWrap(err, "Open file failed")
	}

	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		return debug.ErrorWrap(err, "Copy file failed")
	}

	stat, err := from.Stat()
	if err != nil {
		return debug.ErrorWrap(err, "Stat file failed")
	}

	if err := os.Chmod(dest, stat.Mode()); err != nil {
		return debug.ErrorWrap(err, "Chmod failed")
	}

	return nil
}

//nolint:unused
func copyDir(src, dest string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return debug.ErrorWrap(err, "ReadDir failed")
	}

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return debug.ErrorWrap(err, "Mkdir failed")
	}

	for _, f := range files {
		switch {
		case f.Mode().IsDir():
			if err := copyDir(filepath.Join(src, f.Name()), filepath.Join(dest, f.Name())); err != nil {
				return err
			}

		case f.Mode().IsRegular():
			if err := copyFile(filepath.Join(src, f.Name()), filepath.Join(dest, f.Name())); err != nil {
				return err
			}

		case f.Mode()&os.ModeSymlink != 0:
			actual, err := os.Readlink(filepath.Join(src, f.Name()))

			if err != nil {
				return debug.ErrorWrap(err, "Readlink failed")
			}

			if err := os.Symlink(actual, filepath.Join(dest, f.Name())); err != nil {
				return debug.ErrorWrap(err, "Symlink failed")
			}
		}
	}

	return nil
}
