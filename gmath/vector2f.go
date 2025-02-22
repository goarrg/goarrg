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

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Vector2f[T constraints.Float] struct {
	X, Y T
}

type (
	Vector2f32 = Vector2f[float32]
	Vector2f64 = Vector2f[float64]
)

func Vector2fFromArray[T constraints.Float](a [2]T) Vector2f[T] {
	return Vector2f[T]{
		X: a[0],
		Y: a[1],
	}
}

func (v Vector2f[T]) IsNAN() bool {
	return math.IsNaN(float64(v.X)) || math.IsNaN(float64(v.Y))
}

func (v Vector2f[T]) Abs() Vector2f[T] {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	return v
}

func (v Vector2f[T]) Min(v2 Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: min(v.X, v2.X),
		Y: min(v.Y, v2.Y),
	}
}

func (v Vector2f[T]) Max(v2 Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: max(v.X, v2.X),
		Y: max(v.Y, v2.Y),
	}
}

func (v Vector2f[T]) Clamp(lo, hi Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: max(lo.X, min(v.X, hi.X)),
		Y: max(lo.Y, min(v.Y, hi.Y)),
	}
}

func (v Vector2f[T]) Add(v2 Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v Vector2f[T]) AddUniform(f T) Vector2f[T] {
	return Vector2f[T]{
		X: v.X + f,
		Y: v.Y + f,
	}
}

func (v Vector2f[T]) Subtract(v2 Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v Vector2f[T]) SubtractUniform(f T) Vector2f[T] {
	return Vector2f[T]{
		X: v.X - f,
		Y: v.Y - f,
	}
}

func (v Vector2f[T]) Scale(v2 Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
	}
}

func (v Vector2f[T]) ScaleInverse(v2 Vector2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: v.X / v2.X,
		Y: v.Y / v2.Y,
	}
}

func (v Vector2f[T]) ScaleUniform(s T) Vector2f[T] {
	return Vector2f[T]{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vector2f[T]) ScaleInverseUniform(s T) Vector2f[T] {
	return Vector2f[T]{
		X: v.X / s,
		Y: v.Y / s,
	}
}

func (v Vector2f[T]) Dot(v2 Vector2f[T]) T {
	return (v.X * v2.X) + (v.Y * v2.Y)
}

func (v Vector2f[T]) Magnitude() T {
	return T(math.Sqrt(
		float64(
			(v.X * v.X) + (v.Y * v.Y),
		),
	))
}

func (v Vector2f[T]) Angle(v2 Vector2f[T]) T {
	dot := v.Dot(v2)
	ma := v.Magnitude()
	mb := v2.Magnitude()

	return T(math.Acos(
		float64(dot / (ma * mb)),
	))
}

func (v Vector2f[T]) Normalize() Vector2f[T] {
	m := v.Magnitude()

	if m == 0 || m == 1 {
		return v
	}

	return Vector2f[T]{
		X: v.X / m,
		Y: v.Y / m,
	}
}

func (v Vector2f[T]) ToArray() [2]T {
	return [2]T{v.X, v.Y}
}

func (v Vector2f[T]) ToArrayf32() [2]float32 {
	return [2]float32{float32(v.X), float32(v.Y)}
}
