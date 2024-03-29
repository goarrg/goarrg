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
	"bytes"
	"io"
	"io/fs"
	"math"
	"os"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"goarrg.com/debug"
)

type sys interface {
	bytes() []byte
	uintptr() uintptr
	close()
}

type mmap struct {
	info fileinfo
	sys  sys
	refs *int64
}

var (
	cache  = make(map[string]mmap)
	mtx    = sync.RWMutex{}
	logger = debug.NewLogger("goarrg", "asset")
)

func load(name string) (mmap, error) {
	mtx.RLock()

	if m, ok := cache[name]; ok {
		logger.VPrintf("Loading [%s] from cache", name)
		atomic.AddInt64(m.refs, 1)
		mtx.RUnlock()
		return m, nil
	}

	mtx.RUnlock()
	mtx.Lock()
	defer mtx.Unlock()

	// be 100% sure it wasn't added between the RUnlock() and Lock()
	if m, ok := cache[name]; ok {
		logger.VPrintf("Loading [%s] from cache", name)
		atomic.AddInt64(m.refs, 1)
		return m, nil
	}

	logger.VPrintf("Loading [%s] from disk", name)

	info, err := os.Stat(name)
	if err != nil {
		return mmap{}, debug.ErrorWrapf(err, "Failed to load asset %q", name)
	}

	if info.Size() == 0 {
		return mmap{}, debug.ErrorWrapf(debug.Errorf("Empty file"), "Failed to load asset %q", name)
	}

	if unsafe.Sizeof(int(0)) != unsafe.Sizeof(int64(0)) && info.Size() > math.MaxInt32 {
		return mmap{}, debug.ErrorWrapf(debug.Errorf("File too big"), "Failed to load asset %q", name)
	}

	s, err := mapFile(name, int(info.Size()))
	if err != nil {
		return mmap{}, debug.ErrorWrapf(err, "Failed to load asset %q", name)
	}

	m := mmap{
		info: fileinfo{name: name, size: int(info.Size())},
		sys:  s, refs: new(int64),
	}
	atomic.AddInt64(m.refs, 1)
	cache[name] = m
	return m, nil
}

type fileinfo struct {
	name string
	size int
}

func (i *fileinfo) Name() string {
	return i.name
}

func (i *fileinfo) Size() int64 {
	return int64(i.size)
}

func (i *fileinfo) Mode() fs.FileMode {
	return fs.ModeIrregular
}

func (i *fileinfo) ModTime() time.Time {
	return time.Time{}
}

func (i *fileinfo) IsDir() bool {
	return false
}

func (i *fileinfo) Sys() any {
	return nil
}

type (
	File struct {
		reader *bytes.Reader
		mmap   mmap
	}
)

func (f *File) Name() string {
	return f.mmap.info.name
}

func (f *File) Uintptr() uintptr {
	return f.mmap.sys.uintptr()
}

func (f *File) Size() int {
	return f.mmap.info.size
}

func (f *File) Stat() (fs.FileInfo, error) {
	return &f.mmap.info, nil
}

func (f *File) Len() int {
	return f.reader.Len()
}

func (f *File) Discard(n int64) (discarded int64, err error) {
	return f.Seek(n, io.SeekCurrent)
}

func (f *File) Peek(n int) ([]byte, error) {
	i := f.Size() - f.Len()
	if (i + n) > f.Size() {
		return nil, io.EOF
	}
	return slices.Clone(f.mmap.sys.bytes()[i : i+n]), nil
}

func (f *File) Read(b []byte) (n int, err error) {
	return f.reader.Read(b)
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	return f.reader.ReadAt(b, off)
}

func (f *File) ReadByte() (byte, error) {
	return f.reader.ReadByte()
}

func (f *File) ReadRune() (ch rune, size int, err error) {
	return f.reader.ReadRune()
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *File) UnreadByte() error {
	return f.reader.UnreadByte()
}

func (f *File) UnreadRune() error {
	return f.reader.UnreadRune()
}

func (f *File) WriteTo(w io.Writer) (n int64, err error) {
	return f.reader.WriteTo(w)
}

func (f *File) Close() error {
	if f.mmap.refs == nil {
		return fs.ErrClosed
	}

	if atomic.AddInt64(f.mmap.refs, -1) <= 0 {
		mtx.Lock()
		defer mtx.Unlock()

		if _, ok := cache[f.mmap.info.name]; ok {
			logger.VPrintf("Removing [%s] from cache", f.mmap.info.name)
			delete(cache, f.mmap.info.name)
			f.mmap.sys.close()
		}
	}

	f.mmap.refs = nil
	return nil
}

func Load(name string) (*File, error) {
	m, err := load(name)
	if err != nil {
		return nil, err
	}
	f := File{reader: bytes.NewReader(m.sys.bytes()), mmap: m}
	runtime.SetFinalizer(&f, (*File).Close)

	return &f, nil
}
