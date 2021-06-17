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
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"

	"goarrg.com/debug"
)

type file struct {
	name string
	size int
	fh   windows.Handle
	mh   windows.Handle
	addr uintptr
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

	fh, err := windows.CreateFile(windows.StringToUTF16Ptr(f), windows.GENERIC_READ, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, windows.OPEN_EXISTING, 0x08000000, 0)

	if fh == windows.InvalidHandle {
		return file{}, debug.ErrorWrapf(err, "Failed to map file")
	}

	mh, err := windows.CreateFileMapping(fh, nil, windows.PAGE_READONLY, 0, 0, nil)

	if mh == windows.InvalidHandle {
		if err := windows.CloseHandle(fh); err != nil {
			panic(err)
		}

		return file{}, debug.ErrorWrapf(err, "Failed to map file")
	}

	addr, err := windows.MapViewOfFile(mh, windows.FILE_MAP_READ, 0, 0, 0)

	if addr == 0 {
		if err := windows.CloseHandle(fh); err != nil {
			panic(err)
		}

		if err := windows.CloseHandle(mh); err != nil {
			panic(err)
		}

		return file{}, debug.ErrorWrapf(err, "Failed to map file")
	}

	return file{
		f,
		int(info.Size()),
		fh,
		mh,
		addr,
		*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
			addr, int(info.Size()), int(info.Size()),
		})),
	}, nil
}

func (f *file) close() {
	if f.addr == 0 {
		return
	}
	if err := windows.UnmapViewOfFile(f.addr); err != nil {
		panic(err)
	}
	if err := windows.CloseHandle(f.mh); err != nil {
		panic(err)
	}
	if err := windows.CloseHandle(f.fh); err != nil {
		panic(err)
	}
}
