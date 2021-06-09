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

package qbcl

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"runtime"

	"goarrg.com/asset"
	"goarrg.com/asset/voxel"
	"goarrg.com/debug"
	"goarrg.com/gmath"
)

func init() {
	voxel.RegisterFormat("QB 2", loadQBT)
}

func loadQBT(a asset.Asset) (map[string]voxel.Model, error) {
	var version [2]byte

	r := bufio.NewReader(a.Reader())

	// skip magic
	if _, err := r.Discard(4); err != nil {
		return nil, debug.ErrorWrap(err, "Failed to decode as qbt format")
	}

	if err := binary.Read(r, binary.LittleEndian, version[:]); err != nil {
		return nil, debug.ErrorWrap(err, "Failed to decode as qbt format")
	}

	if version != [2]byte{1, 0} {
		return nil, debug.ErrorWrap(debug.ErrorNew("Either bad or unsupported version %#x", version), "Failed to decode as qbt format")
	}

	if _, err := r.Discard(12); err != nil {
		return nil, debug.ErrorWrap(err, "Failed to decode as qbt format")
	}

	var section [8]byte

	if err := binary.Read(r, binary.LittleEndian, section[:]); err != nil {
		return nil, debug.ErrorWrap(err, "Failed to decode as qbt format")
	}

	if string(section[:]) != "COLORMAP" {
		return nil, debug.ErrorWrap(debug.ErrorNew("Bad file structure"), "Failed to decode as qbt format")
	}

	var colorCount uint32

	if err := binary.Read(r, binary.LittleEndian, &colorCount); err != nil {
		return nil, debug.ErrorWrap(err, "Failed to decode as qbt format")
	}

	if colorCount != 0 {
		return nil, debug.ErrorWrap(debug.ErrorNew("Colormap is unsupported"), "Failed to decode as qbt format")
	}

	if err := binary.Read(r, binary.LittleEndian, section[:]); err != nil {
		return nil, debug.ErrorWrap(err, "Failed to decode as qbt format")
	}

	if string(section[:]) != "DATATREE" {
		return nil, debug.ErrorWrap(debug.ErrorNew("Bad file structure"), "Failed to decode as qbt format")
	}

	m := make(map[string]voxel.Model)

	if err := debug.ErrorWrap(loadNode(r, m), "Failed to decode as qbt format"); err != nil {
		return nil, err
	}

	runtime.KeepAlive(a)

	return m, nil
}

func loadNode(r *bufio.Reader, m map[string]voxel.Model) error {
	var t uint32
	var sz uint32

	if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		return err
	}

	if sz == 0 {
		return debug.ErrorNew("Bad node size %d %d", sz, t)
	}

	switch t {
	case 0:
		return loadMatrix(r, m)
	case 1:
		return loadModel(r, m)
	default:
		return debug.ErrorNew("WTF? %d", t)
	}
}

func loadMatrix(r *bufio.Reader, m map[string]voxel.Model) error {
	var nameSZ uint32
	var name []byte

	if err := binary.Read(r, binary.LittleEndian, &nameSZ); err != nil {
		return err
	}

	if nameSZ == 0 {
		return debug.ErrorWrap(debug.ErrorNew("Bad name size %d", nameSZ), "Failed to load matrix")
	}

	name = make([]byte, nameSZ)

	if _, err := io.ReadFull(r, name); err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	if _, err := r.Discard(12); err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	if _, err := r.Discard(12); err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	if _, err := r.Discard(12); err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	var modelDimensions [3]uint32

	if err := binary.Read(r, binary.LittleEndian, modelDimensions[:]); err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	var dataSZ uint32

	if err := binary.Read(r, binary.LittleEndian, &dataSZ); err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	if dataSZ == 0 {
		return debug.ErrorWrap(debug.ErrorNew("Bad data size %d", dataSZ), "Failed to load matrix")
	}

	data := bytes.NewBuffer(make([]byte, 0, int64(dataSZ)))

	_, err := io.CopyN(data, r, int64(dataSZ))
	if err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	zr, err := zlib.NewReader(data)
	if err != nil {
		return debug.ErrorWrap(err, "Failed to load matrix")
	}

	defer zr.Close()

	voxels := make([]byte, modelDimensions[0]*modelDimensions[1]*modelDimensions[2]*4)

	for x := uint32(0); x < modelDimensions[0]; x++ {
		for z := uint32(0); z < modelDimensions[2]; z++ {
			for y := uint32(0); y < modelDimensions[1]; y++ {
				var voxel [4]byte

				if n, err := zr.Read(voxel[:]); err != nil && n != 4 {
					return debug.ErrorWrap(debug.ErrorWrap(err, "Failed to decode voxels"), "Failed to load matrix %d %d %d", x, y, z)
				}

				if voxel[3] > 0 {
					index := ((modelDimensions[0] - 1 - x) * 4) + ((modelDimensions[0] * y) * 4) + ((modelDimensions[0] * modelDimensions[1] * z) * 4)
					voxel[3] = 255
					_ = append(voxels[index:index:index+4], voxel[:]...)
				}
			}
		}
	}

	m[string(name)] = voxel.Model{
		Data: voxels,
		Size: gmath.Vector3i{
			X: int(modelDimensions[0]),
			Y: int(modelDimensions[1]),
			Z: int(modelDimensions[2]),
		},
	}

	return nil
}

func loadModel(r *bufio.Reader, m map[string]voxel.Model) error {
	var children uint32

	err := binary.Read(r, binary.LittleEndian, &children)
	if err != nil {
		return debug.ErrorWrap(err, "Failed to load model")
	}

	if children == 0 {
		return debug.ErrorWrap(debug.ErrorNew("Invalid children number"), "Failed to load model")
	}

	for i := uint32(0); i < children; i++ {
		err = loadNode(r, m)

		if err != nil {
			return debug.ErrorWrap(err, "Failed to load model")
		}
	}

	return nil
}
