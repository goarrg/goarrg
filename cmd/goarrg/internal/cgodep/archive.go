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

package cgodep

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"goarrg.com/debug"
)

func writeFile(from io.Reader, to string, mode os.FileMode) error {
	f, err := os.OpenFile(to, os.O_CREATE|os.O_RDWR, mode)
	if err != nil {
		return debug.ErrorWrapf(err, "Open file failed")
	}

	defer f.Close()
	_, err = io.Copy(f, from)
	return debug.ErrorWrapf(err, "Copy failed")
}

func extractTARGZ(r io.Reader, dir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return debug.ErrorWrapf(err, "Unknown Error")
	}

	defer gzr.Close()
	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return debug.ErrorWrapf(err, "Unknown Error")
		}

		target := header.Name[strings.IndexAny(header.Name, "/\\")+1:]

		if target == "" {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(dir, target), 0o755); err != nil {
				return debug.ErrorWrapf(err, "Mkdir failed")
			}

		case tar.TypeReg:
			if err := writeFile(tr, filepath.Join(dir, target), os.FileMode(header.Mode)); err != nil {
				return debug.ErrorWrapf(err, "writeFile failed")
			}
		}
	}
}
