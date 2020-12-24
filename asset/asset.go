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
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"goarrg.com/debug"
)

type Config struct {
}

type Asset interface {
	Size() int
	Reader() io.Reader
	Uintptr() uintptr
	Filename() string
}

type asset struct {
	file file
	refs *int64
}

var cache = make(map[string]asset)
var mtx = sync.RWMutex{}

func init() {
	if !(strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], ".test.exe")) {
		if ex, err := os.Executable(); err != nil {
			panic(err)
		} else if err := os.Chdir(filepath.Dir(ex)); err != nil {
			panic(err)
		}
	}
}

func Load(file string) (Asset, error) {
	mtx.RLock()

	if a, ok := cache[file]; ok {
		debug.LogV("Loading asset [%s] from cache", file)
		atomic.AddInt64(a.refs, 1)
		mtx.RUnlock()
		runtime.SetFinalizer(&a, (*asset).close)
		return &a, nil
	}

	mtx.RUnlock()
	mtx.Lock()
	defer mtx.Unlock()

	// be 100% sure it wasn't added between the RUnlock() and Lock()
	if a, ok := cache[file]; ok {
		debug.LogV("Loading asset [%s] from cache", file)
		atomic.AddInt64(a.refs, 1)
		runtime.SetFinalizer(&a, (*asset).close)
		return &a, nil
	}

	debug.LogV("Loading asset [%s] from disk", file)
	f, err := mapFile(file)

	if err != nil {
		return nil, debug.ErrorWrap(err, "Failed to load asset %q", file)
	}

	a := asset{f, new(int64)}
	(*a.refs) = 1
	cache[file] = a
	runtime.SetFinalizer(&a, (*asset).close)

	return &a, nil
}

func (a *asset) Size() int {
	return len(a.file.data)
}

func (a *asset) Reader() io.Reader {
	return bytes.NewReader(a.file.data)
}

func (a *asset) Uintptr() uintptr {
	return uintptr(unsafe.Pointer(&a.file.data[0]))
}

func (a *asset) Filename() string {
	return a.file.name
}

func (a *asset) close() {
	if atomic.AddInt64(a.refs, -1) <= 0 {
		mtx.Lock()
		defer mtx.Unlock()

		if a, ok := cache[a.file.name]; ok {
			debug.LogV("Removing asset [%s] from cache", a.file.name)
			delete(cache, a.file.name)
			a.file.close()
		}
	}
}
