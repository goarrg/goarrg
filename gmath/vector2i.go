/*
Copyright 2025 The goARRG Authors.

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

type Vector2i[T constraints.Integer] struct {
	X, Y T
}

type (
	Vector2int = Vector2i[int]
	Vector2i32 = Vector2i[int32]
	Vector2i64 = Vector2i[int64]
)

func Vector2iFromArray[T constraints.Integer](a [3]T) Vector2i[T] {
	return Vector2i[T]{
		X: a[0],
		Y: a[1],
	}
}

func (v Vector2i[T]) Abs() Vector2i[T] {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	return v
}

func (v Vector2i[T]) Min(v2 Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: min(v.X, v2.X),
		Y: min(v.Y, v2.Y),
	}
}

func (v Vector2i[T]) Max(v2 Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: max(v.X, v2.X),
		Y: max(v.Y, v2.Y),
	}
}

func (v Vector2i[T]) Clamp(lo, hi Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: max(lo.X, min(v.X, hi.X)),
		Y: max(lo.Y, min(v.Y, hi.Y)),
	}
}

func (v Vector2i[T]) Add(v2 Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v Vector2i[T]) AddUniform(i T) Vector2i[T] {
	return Vector2i[T]{
		X: v.X + i,
		Y: v.Y + i,
	}
}

func (v Vector2i[T]) Subtract(v2 Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v Vector2i[T]) SubtractUniform(i T) Vector2i[T] {
	return Vector2i[T]{
		X: v.X - i,
		Y: v.Y - i,
	}
}

func (v Vector2i[T]) Scale(v2 Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
	}
}

func (v Vector2i[T]) ScaleInverse(v2 Vector2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: v.X / v2.X,
		Y: v.Y / v2.Y,
	}
}

func (v Vector2i[T]) ScaleUniform(s T) Vector2i[T] {
	return Vector2i[T]{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vector2i[T]) ScaleInverseUniform(s T) Vector2i[T] {
	return Vector2i[T]{
		X: v.X / s,
		Y: v.Y / s,
	}
}

func (v Vector2i[T]) ToArray() [2]T {
	return [2]T{v.X, v.Y}
}

func (v Vector2i[T]) ToArrayi32() [2]int32 {
	return [2]int32{int32(v.X), int32(v.Y)}
}
