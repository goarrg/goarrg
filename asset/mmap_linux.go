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

package asset

import (
	"math"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"

	"goarrg.com/debug"
)

type file struct {
	name string
	size int
	fd   int
	data []byte
}

func mapFile(f string) (file, error) {
	info, err := os.Stat(f)

	if err != nil {
		return file{}, debug.ErrorWrapf(err, "Failed to map file")
	}

	if info.Size() == 0 {
		return file{}, debug.ErrorWrapf(debug.Errorf("Empty file"), "Failed to map file")
	}

	if unsafe.Sizeof(int(0)) != unsafe.Sizeof(int64(0)) && info.Size() > math.MaxInt32 {
		return file{}, debug.ErrorWrapf(debug.Errorf("File too big"), "Failed to map file")
	}

	fd, err := unix.Open(f, unix.O_RDONLY, 0)

	if err != nil {
		return file{}, debug.ErrorWrapf(err, "Failed to map file")
	}

	data, err := unix.Mmap(fd, 0, int(info.Size()), unix.PROT_READ, unix.MAP_SHARED)

	if err != nil {
		unix.Close(fd)
		return file{}, debug.ErrorWrapf(err, "Failed to map file")
	}

	return file{
		f,
		int(info.Size()),
		fd,
		data,
	}, nil
}

func (f *file) close() {
	debug.LogErr(debug.ErrorWrapf(unix.Munmap(f.data), "Failed to unmap file"))
	unix.Close(f.fd)
}
