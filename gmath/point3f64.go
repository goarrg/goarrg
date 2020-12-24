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

package gmath

import (
	"math"
)

type Point3f64 struct {
	X, Y, Z float64
}

func Point3f64FromArray(a [3]float64) Point3f64 {
	return Point3f64{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (p Point3f64) VectorTo(p2 Point3f64) Vector3f64 {
	return Vector3f64{
		X: p2.X - p.X,
		Y: p2.Y - p.Y,
		Z: p2.Z - p.Z,
	}
}

func (p Point3f64) Translate(v Vector3f64) Point3f64 {
	return Point3f64{
		X: p.X + v.X,
		Y: p.Y + v.Y,
		Z: p.Z + v.Z,
	}
}

func (p Point3f64) IsNAN() bool {
	return math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsNaN(p.Z)
}

func (p Point3f64) ToArray() [3]float64 {
	return [3]float64{p.X, p.Y, p.Z}
}

func (p Point3f64) ToArrayf32() [3]float32 {
	return [3]float32{float32(p.X), float32(p.Y), float32(p.Z)}
}
