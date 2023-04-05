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
	"unsafe"

	"golang.org/x/sys/unix"

	"goarrg.com/debug"
)

type sysLinux struct {
	fd   int
	addr uintptr
	data []byte
}

func mapFile(name string, size int) (sys, error) {
	fd, err := unix.Open(name, unix.O_RDONLY, 0)
	if err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to map file")
	}

	data, err := unix.Mmap(fd, 0, size, unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		unix.Close(fd)
		return nil, debug.ErrorWrapf(err, "Failed to map file")
	}

	return &sysLinux{
		fd,
		uintptr(unsafe.Pointer(unsafe.SliceData(data))),
		data,
	}, nil
}

func (f *sysLinux) bytes() []byte {
	return f.data
}

func (f *sysLinux) uintptr() uintptr {
	return f.addr
}

func (f *sysLinux) close() {
	if err := debug.ErrorWrapf(unix.Munmap(f.data), "Failed to unmap file"); err != nil {
		panic(err)
	}
	unix.Close(f.fd)
}
