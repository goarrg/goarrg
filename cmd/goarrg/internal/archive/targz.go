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

package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"strings"

	"goarrg.com/debug"
)

func extractTARGZHere(r io.Reader) error {
	gzr, err := gzip.NewReader(r)

	if err != nil {
		return debug.ErrorWrap(err, "Unknown Error")
	}

	defer gzr.Close()
	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return debug.ErrorWrap(err, "Unknown Error")
		}

		target := header.Name[strings.IndexAny(header.Name, "/\\")+1:]

		if target == "" {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return debug.ErrorWrap(err, "Mkdir failed")
			}

		case tar.TypeReg:
			if err := copyFile(tr, target, os.FileMode(header.Mode)); err != nil {
				return debug.ErrorWrap(err, "copyFile failed")
			}

		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return debug.ErrorWrap(err, "Symlink failed")
			}
		}
	}
}
