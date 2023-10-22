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

	"golang.org/x/sys/windows"

	"goarrg.com/debug"
)

type sysWindows struct {
	fh   windows.Handle
	mh   windows.Handle
	addr uintptr
	data []byte
}

func mapFile(name string, size int) (sys, error) {
	var err error
	s := sysWindows{fh: windows.InvalidHandle, mh: windows.InvalidHandle}

	s.fh, err = windows.CreateFile(windows.StringToUTF16Ptr(name), windows.GENERIC_READ, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil,
		windows.OPEN_EXISTING, windows.FILE_FLAG_SEQUENTIAL_SCAN, 0)
	if s.fh == windows.InvalidHandle {
		return nil, debug.ErrorWrapf(err, "Failed to map file")
	}

	s.mh, err = windows.CreateFileMapping(s.fh, nil, windows.PAGE_READONLY, 0, 0, nil)
	if s.mh == windows.InvalidHandle {
		s.close()
		return nil, debug.ErrorWrapf(err, "Failed to map file")
	}

	s.addr, err = windows.MapViewOfFile(s.mh, windows.FILE_MAP_READ, 0, 0, 0)
	if s.addr == 0 {
		s.close()
		return nil, debug.ErrorWrapf(err, "Failed to map file")
	}

	s.data = unsafe.Slice(
		(*byte)(unsafe.Pointer(s.addr)), size,
	)
	return &s, nil
}

func (f *sysWindows) bytes() []byte {
	return f.data
}

func (f *sysWindows) uintptr() uintptr {
	return f.addr
}

func (f *sysWindows) close() {
	if f.addr != 0 {
		if err := windows.UnmapViewOfFile(f.addr); err != nil {
			panic(err)
		}
	}

	if f.mh != windows.InvalidHandle {
		if err := windows.CloseHandle(f.mh); err != nil {
			panic(err)
		}
	}

	if f.fh != windows.InvalidHandle {
		if err := windows.CloseHandle(f.fh); err != nil {
			panic(err)
		}
	}
}
