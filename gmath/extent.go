/*
Copyright 2024 The goARRG Authors.

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

type Extent2f[T constraints.Float] struct {
	X, Y T
}

type (
	Extent2f32 = Extent2f[float32]
	Extent2f64 = Extent2f[float64]
)

func (e Extent2f[T]) InRange(min, max Extent2f[T]) bool {
	return InRange(e.X, min.X, max.X) &&
		InRange(e.Y, min.Y, max.Y)
}

func (e Extent2f[T]) Area() T {
	return e.X * e.Y
}

func (e Extent2f[T]) Min(o Extent2f[T]) Extent2f[T] {
	return Extent2f[T]{
		X: min(e.X, o.X),
		Y: min(e.Y, o.Y),
	}
}

func (e Extent2f[T]) Max(o Extent2f[T]) Extent2f[T] {
	return Extent2f[T]{
		X: max(e.X, o.X),
		Y: max(e.Y, o.Y),
	}
}

func (e Extent2f[T]) Clamp(lo, hi Extent2f[T]) Extent2f[T] {
	return Extent2f[T]{
		X: Clamp(e.X, lo.X, hi.X),
		Y: Clamp(e.Y, lo.Y, hi.Y),
	}
}

func (b Extent2f[T]) CheckPoint(p Point2f[T]) bool {
	return InRange(p.X, 0, b.X) &&
		InRange(p.Y, 0, b.Y)
}

func (b Extent2f[T]) CheckVector(v Vector2f[T]) bool {
	return InRange(v.X, 0, b.X) &&
		InRange(v.Y, 0, b.Y)
}

func (b Extent2f[T]) ClampPoint(p Point2f[T]) Point2f[T] {
	return Point2f[T]{
		X: Clamp(p.X, 0, b.X),
		Y: Clamp(p.Y, 0, b.Y),
	}
}

func (b Extent2f[T]) ClampVector(v Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: Clamp(v.X, 0, b.X),
		Y: Clamp(v.Y, 0, b.Y),
	}
}

type Extent2i[T constraints.Integer] struct {
	X, Y T
}

type (
	Extent2int = Extent2i[int]
	Extent2i32 = Extent2i[int32]
	Extent2i64 = Extent2i[int64]

	Extent2uint = Extent2i[uint]
	Extent2u32  = Extent2i[uint32]
	Extent2u64  = Extent2i[uint64]
)

func (e Extent2i[T]) InRange(min, max Extent2i[T]) bool {
	return InRange(e.X, min.X, max.X) &&
		InRange(e.Y, min.Y, max.Y)
}

func (e Extent2i[T]) Area() T {
	return e.X * e.Y
}

func (e Extent2i[T]) Min(o Extent2i[T]) Extent2i[T] {
	return Extent2i[T]{
		X: min(e.X, o.X),
		Y: min(e.Y, o.Y),
	}
}

func (e Extent2i[T]) Max(o Extent2i[T]) Extent2i[T] {
	return Extent2i[T]{
		X: max(e.X, o.X),
		Y: max(e.Y, o.Y),
	}
}

func (e Extent2i[T]) Clamp(lo, hi Extent2i[T]) Extent2i[T] {
	return Extent2i[T]{
		X: Clamp(e.X, lo.X, hi.X),
		Y: Clamp(e.Y, lo.Y, hi.Y),
	}
}

func (b Extent2i[T]) CheckPoint(p Point2i[T]) bool {
	return InRange(p.X, 0, b.X) &&
		InRange(p.Y, 0, b.Y)
}

func (b Extent2i[T]) CheckVector(v Vector2i[T]) bool {
	return InRange(v.X, 0, b.X) &&
		InRange(v.Y, 0, b.Y)
}

func (b Extent2i[T]) ClampPoint(p Point2i[T]) Point2i[T] {
	return Point2i[T]{
		X: Clamp(p.X, 0, b.X),
		Y: Clamp(p.Y, 0, b.Y),
	}
}

func (b Extent2i[T]) ClampVector(v Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: Clamp(v.X, 0, b.X),
		Y: Clamp(v.Y, 0, b.Y),
	}
}

type Extent3f[T constraints.Float] struct {
	X, Y, Z T
}

type (
	Extent3f32 = Extent3f[float32]
	Extent3f64 = Extent3f[float64]
)

func (e Extent3f[T]) InRange(min, max Extent3f[T]) bool {
	return InRange(e.X, min.X, max.X) &&
		InRange(e.Y, min.Y, max.Y) &&
		InRange(e.Z, min.Z, max.Z)
}

func (e Extent3f[T]) Volume() T {
	return e.X * e.Y * e.Z
}

func (e Extent3f[T]) Min(o Extent3f[T]) Extent3f[T] {
	return Extent3f[T]{
		X: min(e.X, o.X),
		Y: min(e.Y, o.Y),
		Z: min(e.Z, o.Z),
	}
}

func (e Extent3f[T]) Max(o Extent3f[T]) Extent3f[T] {
	return Extent3f[T]{
		X: max(e.X, o.X),
		Y: max(e.Y, o.Y),
		Z: max(e.Z, o.Z),
	}
}

func (e Extent3f[T]) Clamp(lo, hi Extent3f[T]) Extent3f[T] {
	return Extent3f[T]{
		X: Clamp(e.X, lo.X, hi.X),
		Y: Clamp(e.Y, lo.Y, hi.Y),
		Z: Clamp(e.Z, lo.Z, hi.Z),
	}
}

func (b Extent3f[T]) CheckPoint(p Point3f[T]) bool {
	return InRange(p.X, 0, b.X) &&
		InRange(p.Y, 0, b.Y) &&
		InRange(p.Z, 0, b.Z)
}

func (b Extent3f[T]) CheckVector(v Vector3f[T]) bool {
	return InRange(v.X, 0, b.X) &&
		InRange(v.Y, 0, b.Y) &&
		InRange(v.Z, 0, b.Z)
}

func (b Extent3f[T]) ClampPoint(p Point3f[T]) Point3f[T] {
	return Point3f[T]{
		X: Clamp(p.X, 0, b.X),
		Y: Clamp(p.Y, 0, b.Y),
		Z: Clamp(p.Z, 0, b.Z),
	}
}

func (b Extent3f[T]) ClampVector(v Vector3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: Clamp(v.X, 0, b.X),
		Y: Clamp(v.Y, 0, b.Y),
		Z: Clamp(v.Z, 0, b.Z),
	}
}

type Extent3i[T constraints.Integer] struct {
	X, Y, Z T
}

type (
	Extent3int = Extent3i[int]
	Extent3i32 = Extent3i[int32]
	Extent3i64 = Extent3i[int64]

	Extent3uint = Extent3i[uint]
	Extent3u32  = Extent3i[uint32]
	Extent3u64  = Extent3i[uint64]
)

func (e Extent3i[T]) InRange(min, max Extent3i[T]) bool {
	return InRange(e.X, min.X, max.X) &&
		InRange(e.Y, min.Y, max.Y) &&
		InRange(e.Z, min.Z, max.Z)
}

func (e Extent3i[T]) Volume() T {
	return e.X * e.Y * e.Z
}

func (e Extent3i[T]) Min(o Extent3i[T]) Extent3i[T] {
	return Extent3i[T]{
		X: min(e.X, o.X),
		Y: min(e.Y, o.Y),
		Z: min(e.Z, o.Z),
	}
}

func (e Extent3i[T]) Max(o Extent3i[T]) Extent3i[T] {
	return Extent3i[T]{
		X: max(e.X, o.X),
		Y: max(e.Y, o.Y),
		Z: max(e.Z, o.Z),
	}
}

func (e Extent3i[T]) Clamp(lo, hi Extent3i[T]) Extent3i[T] {
	return Extent3i[T]{
		X: Clamp(e.X, lo.X, hi.X),
		Y: Clamp(e.Y, lo.Y, hi.Y),
		Z: Clamp(e.Z, lo.Z, hi.Z),
	}
}

func (b Extent3i[T]) CheckPoint(p Point3i[T]) bool {
	return InRange(p.X, 0, b.X) &&
		InRange(p.Y, 0, b.Y) &&
		InRange(p.Z, 0, b.Z)
}

func (b Extent3i[T]) CheckVector(v Vector3i[T]) bool {
	return InRange(v.X, 0, b.X) &&
		InRange(v.Y, 0, b.Y) &&
		InRange(v.Z, 0, b.Z)
}

func (b Extent3i[T]) ClampPoint(p Point3i[T]) Point3i[T] {
	return Point3i[T]{
		X: Clamp(p.X, 0, b.X),
		Y: Clamp(p.Y, 0, b.Y),
		Z: Clamp(p.Z, 0, b.Z),
	}
}

func (b Extent3i[T]) ClampVector(v Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: Clamp(v.X, 0, b.X),
		Y: Clamp(v.Y, 0, b.Y),
		Z: Clamp(v.Z, 0, b.Z),
	}
}
