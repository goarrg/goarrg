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

type Bounds[T constraints.Float | constraints.Integer] struct {
	Min T
	Max T
}

func (b Bounds[T]) CheckValue(v T) bool {
	return InRange(v, b.Min, b.Max)
}

type Bounds3f[T constraints.Float] struct {
	Min Vector3f[T]
	Max Vector3f[T]
}

func (b Bounds3f[T]) CheckPoint(p Point3f[T]) bool {
	return InRange(p.X, b.Min.X, b.Max.X) &&
		InRange(p.Y, b.Min.Y, b.Max.Y) &&
		InRange(p.Z, b.Min.Z, b.Max.Z)
}

func (b Bounds3f[T]) CheckVector(v Vector3f[T]) bool {
	return InRange(v.X, b.Min.X, b.Max.X) &&
		InRange(v.Y, b.Min.Y, b.Max.Y) &&
		InRange(v.Z, b.Min.Z, b.Max.Z)
}

type (
	Bounds3f32 = Bounds3f[float32]
	Bounds3f64 = Bounds3f[float64]
)

type Bounds3i[T constraints.Integer] struct {
	Min Vector3i[T]
	Max Vector3i[T]
}

func (b Bounds3i[T]) CheckPoint(p Point3i[T]) bool {
	return InRange(p.X, b.Min.X, b.Max.X) &&
		InRange(p.Y, b.Min.Y, b.Max.Y) &&
		InRange(p.Z, b.Min.Z, b.Max.Z)
}

func (b Bounds3i[T]) CheckVector(v Vector3i[T]) bool {
	return InRange(v.X, b.Min.X, b.Max.X) &&
		InRange(v.Y, b.Min.Y, b.Max.Y) &&
		InRange(v.Z, b.Min.Z, b.Max.Z)
}

type (
	Bounds3int = Bounds3i[int]
	Bounds3i32 = Bounds3i[int32]
	Bounds3i64 = Bounds3i[int64]
)
