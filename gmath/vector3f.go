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

type Vector3f[T constraints.Float] struct {
	X, Y, Z T
}

type (
	Vector3f32 = Vector3f[float32]
	Vector3f64 = Vector3f[float64]
)

func Vector3FromArray[T constraints.Float](a [3]T) Vector3f[T] {
	return Vector3f[T]{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (v Vector3f[T]) IsNAN() bool {
	return math.IsNaN(float64(v.X)) || math.IsNaN(float64(v.Y)) || math.IsNaN(float64(v.Z))
}

func (v Vector3f[T]) Add(v2 Vector3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vector3f[T]) AddUniform(f T) Vector3f[T] {
	return Vector3f[T]{
		X: v.X + f,
		Y: v.Y + f,
		Z: v.Z + f,
	}
}

func (v Vector3f[T]) Subtract(v2 Vector3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

func (v Vector3f[T]) SubtractUniform(f T) Vector3f[T] {
	return Vector3f[T]{
		X: v.X - f,
		Y: v.Y - f,
		Z: v.Z - f,
	}
}

func (v Vector3f[T]) Scale(v2 Vector3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
		Z: v.Z * v2.Z,
	}
}

func (v Vector3f[T]) ScaleInverse(v2 Vector3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: v.X / v2.X,
		Y: v.Y / v2.Y,
		Z: v.Z / v2.Z,
	}
}

func (v Vector3f[T]) ScaleUniform(s T) Vector3f[T] {
	return Vector3f[T]{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func (v Vector3f[T]) ScaleInverseUniform(s T) Vector3f[T] {
	return Vector3f[T]{
		X: v.X / s,
		Y: v.Y / s,
		Z: v.Z / s,
	}
}

func (v Vector3f[T]) Dot(v2 Vector3f[T]) T {
	return (v.X * v2.X) + (v.Y * v2.Y) + (v.Z * v2.Z)
}

func (v Vector3f[T]) Magnitude() T {
	return T(math.Sqrt(
		float64(
			(v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z),
		),
	))
}

func (v Vector3f[T]) Angle(v2 Vector3f[T]) T {
	dot := v.Dot(v2)
	ma := v.Magnitude()
	mb := v2.Magnitude()

	return T(math.Acos(
		float64(dot / (ma * mb)),
	))
}

func (v Vector3f[T]) AngleAxis(v2 Vector3f[T]) (T, Vector3f[T]) {
	return v.Angle(v2), v.Normalize().Cross(v2.Normalize()).Normalize()
}

func (v Vector3f[T]) Normalize() Vector3f[T] {
	m := v.Magnitude()

	if m == 0 || m == 1 {
		return v
	}

	return Vector3f[T]{
		X: v.X / m,
		Y: v.Y / m,
		Z: v.Z / m,
	}
}

func (v Vector3f[T]) Cross(v2 Vector3f[T]) Vector3f[T] {
	return Vector3f[T]{
		X: (v.Y * v2.Z) - (v.Z * v2.Y),
		Y: (v.Z * v2.X) - (v.X * v2.Z),
		Z: (v.X * v2.Y) - (v.Y * v2.X),
	}
}

func (v Vector3f[T]) OrthoNormalize(tangent Vector3f[T]) Vector3f[T] {
	v = v.Normalize()
	tangent = tangent.Normalize()

	return tangent.Subtract(v.ScaleUniform(tangent.Dot(v))).Normalize()
}

func (v Vector3f[T]) ToArray() [3]T {
	return [3]T{v.X, v.Y, v.Z}
}

func (v Vector3f[T]) ToArrayf32() [3]float32 {
	return [3]float32{float32(v.X), float32(v.Y), float32(v.Z)}
}
