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

type Vector3i[T constraints.Integer] struct {
	X, Y, Z T
}

type (
	Vector3int = Vector3i[int]
	Vector3i32 = Vector3i[int32]
	Vector3i64 = Vector3i[int64]
)

func Vector3iFromArray[T constraints.Integer](a [3]T) Vector3i[T] {
	return Vector3i[T]{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (v Vector3i[T]) Abs() Vector3i[T] {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	if v.Z < 0 {
		v.Z = -v.Z
	}
	return v
}

func (v Vector3i[T]) Min(v2 Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: min(v.X, v2.X),
		Y: min(v.Y, v2.Y),
		Z: min(v.Z, v2.Z),
	}
}

func (v Vector3i[T]) Max(v2 Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: max(v.X, v2.X),
		Y: max(v.Y, v2.Y),
		Z: max(v.Z, v2.Z),
	}
}

func (v Vector3i[T]) Clamp(lo, hi Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: max(lo.X, min(v.X, hi.X)),
		Y: max(lo.Y, min(v.Y, hi.Y)),
		Z: max(lo.Z, min(v.Z, hi.Z)),
	}
}

func (v Vector3i[T]) Add(v2 Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vector3i[T]) AddUniform(i T) Vector3i[T] {
	return Vector3i[T]{
		X: v.X + i,
		Y: v.Y + i,
		Z: v.Z + i,
	}
}

func (v Vector3i[T]) Subtract(v2 Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

func (v Vector3i[T]) SubtractUniform(i T) Vector3i[T] {
	return Vector3i[T]{
		X: v.X - i,
		Y: v.Y - i,
		Z: v.Z - i,
	}
}

func (v Vector3i[T]) Scale(v2 Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
		Z: v.Z * v2.Z,
	}
}

func (v Vector3i[T]) ScaleInverse(v2 Vector3i[T]) Vector3i[T] {
	return Vector3i[T]{
		X: v.X / v2.X,
		Y: v.Y / v2.Y,
		Z: v.Z / v2.Z,
	}
}

func (v Vector3i[T]) ScaleUniform(s T) Vector3i[T] {
	return Vector3i[T]{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func (v Vector3i[T]) ScaleInverseUniform(s T) Vector3i[T] {
	return Vector3i[T]{
		X: v.X / s,
		Y: v.Y / s,
		Z: v.Z / s,
	}
}

func (v Vector3i[T]) ToArray() [3]T {
	return [3]T{v.X, v.Y, v.Z}
}

func (v Vector3i[T]) ToArrayi32() [3]int32 {
	return [3]int32{int32(v.X), int32(v.Y), int32(v.Z)}
}
