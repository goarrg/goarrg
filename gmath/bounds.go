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

func (b Bounds[T]) ClampValue(v T) T {
	return Clamp(v, b.Min, b.Max)
}

type Bounds2f[T constraints.Float] struct {
	Min Vector2f[T]
	Max Vector2f[T]
}

func (b Bounds2f[T]) CheckPoint(p Point2f[T]) bool {
	return InRange(p.X, b.Min.X, b.Max.X) &&
		InRange(p.Y, b.Min.Y, b.Max.Y)
}

func (b Bounds2f[T]) CheckVector(v Vector2f[T]) bool {
	return InRange(v.X, b.Min.X, b.Max.X) &&
		InRange(v.Y, b.Min.Y, b.Max.Y)
}

func (b Bounds2f[T]) ClampPoint(p Point2f[T]) Point2f[T] {
	return Point2f[T]{
		X: Clamp(p.X, b.Min.X, b.Max.X),
		Y: Clamp(p.Y, b.Min.Y, b.Max.Y),
	}
}

func (b Bounds2f[T]) ClampVector(v Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: Clamp(v.X, b.Min.X, b.Max.X),
		Y: Clamp(v.Y, b.Min.Y, b.Max.Y),
	}
}

type (
	Bounds2f32 = Bounds2f[float32]
	Bounds2f64 = Bounds2f[float64]
)

type Bounds2i[T constraints.Integer] struct {
	Min Vector2i[T]
	Max Vector2i[T]
}

func (b Bounds2i[T]) CheckPoint(p Point2i[T]) bool {
	return InRange(p.X, b.Min.X, b.Max.X) &&
		InRange(p.Y, b.Min.Y, b.Max.Y)
}

func (b Bounds2i[T]) CheckVector(v Vector2i[T]) bool {
	return InRange(v.X, b.Min.X, b.Max.X) &&
		InRange(v.Y, b.Min.Y, b.Max.Y)
}

func (b Bounds2i[T]) ClampPoint(p Point2i[T]) Point2i[T] {
	return Point2i[T]{
		X: Clamp(p.X, b.Min.X, b.Max.X),
		Y: Clamp(p.Y, b.Min.Y, b.Max.Y),
	}
}

func (b Bounds2i[T]) ClampVector(v Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: Clamp(v.X, b.Min.X, b.Max.X),
		Y: Clamp(v.Y, b.Min.Y, b.Max.Y),
	}
}

type (
	Bounds2int = Bounds2i[int]
	Bounds2i32 = Bounds2i[int32]
	Bounds2i64 = Bounds2i[int64]
)

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

func (b Bounds3f[T]) ClampPoint(p Point3f[T]) Point3f[T] {
	return Point3f[T]{
		X: Clamp(p.X, b.Min.X, b.Max.X),
		Y: Clamp(p.Y, b.Min.Y, b.Max.Y),
		Z: Clamp(p.Z, b.Min.Z, b.Max.Z),
	}
}

func (b Bounds3f[T]) ClampVector(v Vector3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: Clamp(v.X, b.Min.X, b.Max.X),
		Y: Clamp(v.Y, b.Min.Y, b.Max.Y),
		Z: Clamp(v.Z, b.Min.Z, b.Max.Z),
	}
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

func (b Bounds3i[T]) ClampPoint(p Point3i[T]) Point3i[T] {
	return Point3i[T]{
		X: Clamp(p.X, b.Min.X, b.Max.X),
		Y: Clamp(p.Y, b.Min.Y, b.Max.Y),
		Z: Clamp(p.Z, b.Min.Z, b.Max.Z),
	}
}

func (b Bounds3i[T]) ClampVector(v Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: Clamp(v.X, b.Min.X, b.Max.X),
		Y: Clamp(v.Y, b.Min.Y, b.Max.Y),
		Z: Clamp(v.Z, b.Min.Z, b.Max.Z),
	}
}

type (
	Bounds3int = Bounds3i[int]
	Bounds3i32 = Bounds3i[int32]
	Bounds3i64 = Bounds3i[int64]
)
