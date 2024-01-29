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

import "golang.org/x/exp/constraints"

type Point3i[T constraints.Integer] struct {
	X, Y, Z T
}

type (
	Point3int = Point3i[int]
	Point3i32 = Point3i[int32]
	Point3i64 = Point3i[int64]
)

func Point3iFromArray[T constraints.Integer](a [3]T) Point3i[T] {
	return Point3i[T]{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (p Point3i[T]) VectorTo(p2 Point3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: p2.X - p.X,
		Y: p2.Y - p.Y,
		Z: p2.Z - p.Z,
	}
}

func (p Point3i[T]) Add(v Vector3i[T]) Point3i[T] {
	return Point3i[T]{
		X: p.X + v.X,
		Y: p.Y + v.Y,
		Z: p.Z + v.Z,
	}
}

func (p Point3i[T]) Subtract(v Vector3i[T]) Point3i[T] {
	return Point3i[T]{
		X: p.X - v.X,
		Y: p.Y - v.Y,
		Z: p.Z - v.Z,
	}
}

func (p Point3i[T]) ToArray() [3]T {
	return [3]T{p.X, p.Y, p.Z}
}
