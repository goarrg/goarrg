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

package voxel

import (
	"bufio"
	"sync"
	"sync/atomic"

	"goarrg.com/asset"
	"goarrg.com/debug"
	"goarrg.com/gmath"
)

type Model struct {
	Data []byte
	Size gmath.Vector3i
}

type Collection struct {
	Models map[string]Model
}

type format struct {
	magic  []byte
	decode func(asset.Asset) (map[string]Model, error)
}

var (
	mtx     = sync.Mutex{}
	formats = atomic.Value{}
)

func RegisterFormat(magic string, decode func(asset.Asset) (map[string]Model, error)) {
	mtx.Lock()
	f, _ := formats.Load().([]format)
	formats.Store(append(f, format{[]byte(magic), decode}))
	mtx.Unlock()
}

func Load(file string) (*Collection, error) {
	a, err := asset.Load(file)
	if err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to load voxel collection")
	}

	r := bufio.NewReader(a.Reader())
	formats, _ := formats.Load().([]format)

formats:
	for _, f := range formats {
		if a.Size() < len(f.magic) {
			continue
		}

		magic, err := r.Peek(len(f.magic))
		if err != nil {
			return nil, debug.ErrorWrapf(err, "Failed to load voxel collection")
		}

		for i, b := range f.magic {
			if b != magic[i] && b != '?' {
				continue formats
			}
		}

		models, err := f.decode(a)
		if err != nil {
			return nil, debug.ErrorWrapf(err, "Failed to load voxel collection")
		}

		return &Collection{
			models,
		}, nil
	}

	return nil, debug.Errorf("Failed to load voxel collection, unknown format")
}
