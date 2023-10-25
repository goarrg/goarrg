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

	"golang.org/x/exp/constraints"
)

type Point3f[T constraints.Float] struct {
	X, Y, Z T
}

type (
	Point3f32 = Point3f[float32]
	Point3f64 = Point3f[float64]
)

func Point3fFromArray[T constraints.Float](a [3]T) Point3f[T] {
	return Point3f[T]{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (p Point3f[T]) VectorTo(p2 Point3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: p2.X - p.X,
		Y: p2.Y - p.Y,
		Z: p2.Z - p.Z,
	}
}

func (p Point3f[T]) Translate(v Vector3f[T]) Point3f[T] {
	return Point3f[T]{
		X: p.X + v.X,
		Y: p.Y + v.Y,
		Z: p.Z + v.Z,
	}
}

func (p Point3f[T]) IsNAN() bool {
	return math.IsNaN(float64(p.X)) || math.IsNaN(float64(p.Y)) || math.IsNaN(float64(p.Z))
}

func (p Point3f[T]) ToArray() [3]T {
	return [3]T{p.X, p.Y, p.Z}
}

func (p Point3f[T]) ToArrayf32() [3]float32 {
	return [3]float32{float32(p.X), float32(p.Y), float32(p.Z)}
}
