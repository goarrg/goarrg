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

type Transform[T constraints.Float] struct {
	Pos   Point3f[T]
	Rot   Quaternion[T]
	Scale Vector3f[T]
	Pivot Point3f[T]
}

type (
	Transformf32 = Transform[float32]
	Transformf64 = Transform[float64]
)

func (t *Transform[T]) TransformPoint(p Point3f[T]) Point3f[T] {
	return t.Pos.Add(t.Rot.Rotate(Vector3f[T](p.Subtract(Vector3f[T](t.Pivot))).Scale(t.Scale)))
}

func (t *Transform[T]) TransformVector(v Vector3f[T]) Vector3f[T] {
	return t.Rot.Rotate(v.Scale(t.Scale))
}

func (t *Transform[T]) TransformDirection(v Vector3f[T]) Vector3f[T] {
	return t.Rot.Rotate(v)
}

func (t *Transform[T]) LookAt(up Vector3f[T], target Point3f[T]) {
	m2 := t.Pos.VectorTo(target).Normalize()
	m0 := up.Cross(m2).Normalize()
	m1 := m2.Cross(m0)
	if trace := m0.X + m1.Y + m2.Z; trace > 0 {
		s := T(math.Sqrt(float64(trace + 1)))
		t.Rot.W = s * 0.5
		s = 0.5 / s
		t.Rot.X = (m1.Z - m2.Y) * s
		t.Rot.Y = (m2.X - m0.Z) * s
		t.Rot.Z = (m0.Y - m1.X) * s
		return
	}
	if (m0.X >= m1.Y) && (m0.X >= m2.Z) {
		s := T(math.Sqrt(float64(1 + m0.X - m1.Y - m2.Z)))
		t.Rot.X = 0.5 * s
		s = 0.5 / s
		t.Rot.Y = (m0.Y + m1.X) * s
		t.Rot.Z = (m0.Z + m2.X) * s
		t.Rot.W = (m1.Z - m2.Y) * s
		return
	}
	if m1.Y > m2.Z {
		s := T(math.Sqrt(float64(1 + m1.Y - m0.X - m2.Z)))
		t.Rot.Y = 0.5 * s
		s = 0.5 / s
		t.Rot.X = (m1.X + m0.Y) * s
		t.Rot.Z = (m2.Y + m1.Z) * s
		t.Rot.W = (m2.X - m0.Z) * s
		return
	}
	s := T(math.Sqrt(float64(1 + m2.Z - m0.X - m1.Y)))
	t.Rot.Z = 0.5 * s
	s = 0.5 / s
	t.Rot.X = (m2.X + m0.Z) * s
	t.Rot.Y = (m2.Y + m1.Z) * s
	t.Rot.W = (m0.Y - m1.X) * s
}

func (t *Transform[T]) PivotMatrix() Matrix4x4f[T] {
	return Matrix4x4f[T]{
		{1, 0, 0, -t.Pivot.X},
		{0, 1, 0, -t.Pivot.Y},
		{0, 0, 1, -t.Pivot.Z},
		{0, 0, 0, 1},
	}
}

func (t *Transform[T]) TranslationMatrix() Matrix4x4f[T] {
	return Matrix4x4f[T]{
		{1, 0, 0, t.Pos.X},
		{0, 1, 0, t.Pos.Y},
		{0, 0, 1, t.Pos.Z},
		{0, 0, 0, 1},
	}
}

func (t *Transform[T]) RotationMatrix() Matrix4x4f[T] {
	x2 := t.Rot.X * t.Rot.X * 2
	y2 := t.Rot.Y * t.Rot.Y * 2
	z2 := t.Rot.Z * t.Rot.Z * 2

	xy := t.Rot.X * t.Rot.Y * 2
	wz := t.Rot.W * t.Rot.Z * 2
	xz := t.Rot.X * t.Rot.Z * 2
	wy := t.Rot.W * t.Rot.Y * 2
	yz := t.Rot.Y * t.Rot.Z * 2
	wx := t.Rot.W * t.Rot.X * 2

	return Matrix4x4f[T]{
		{1 - y2 - z2, xy - wz, xz + wy, 0},
		{xy + wz, 1 - x2 - z2, yz - wx, 0},
		{xz - wy, yz + wx, 1 - x2 - y2, 0},
		{0, 0, 0, 1},
	}
}

func (t *Transform[T]) ScaleMatrix() Matrix4x4f[T] {
	return Matrix4x4f[T]{
		{t.Scale.X, 0, 0, 0},
		{0, t.Scale.Y, 0, 0},
		{0, 0, t.Scale.Z, 0},
		{0, 0, 0, 1},
	}
}

func (t *Transform[T]) ModelMatrix() Matrix4x4f[T] {
	x2 := t.Rot.X * t.Rot.X * 2
	y2 := t.Rot.Y * t.Rot.Y * 2
	z2 := t.Rot.Z * t.Rot.Z * 2

	xy := t.Rot.X * t.Rot.Y * 2
	wz := t.Rot.W * t.Rot.Z * 2
	xz := t.Rot.X * t.Rot.Z * 2
	wy := t.Rot.W * t.Rot.Y * 2
	yz := t.Rot.Y * t.Rot.Z * 2
	wx := t.Rot.W * t.Rot.X * 2

	m0 := Vector3f[T]{
		X: 1 - y2 - z2, Y: xy - wz, Z: xz + wy,
	}.Scale(t.Scale)
	m1 := Vector3f[T]{
		X: xy + wz, Y: 1 - x2 - z2, Z: yz - wx,
	}.Scale(t.Scale)
	m2 := Vector3f[T]{
		X: xz - wy, Y: yz + wx, Z: 1 - x2 - y2,
	}.Scale(t.Scale)

	return Matrix4x4f[T]{
		{m0.X, m0.Y, m0.Z, t.Pos.X - m0.Dot(Vector3f[T](t.Pivot))},
		{m1.X, m1.Y, m1.Z, t.Pos.Y - m1.Dot(Vector3f[T](t.Pivot))},
		{m2.X, m2.Y, m2.Z, t.Pos.Z - m2.Dot(Vector3f[T](t.Pivot))},
		{0, 0, 0, 1},
	}
}

func (t *Transform[T]) ModelInverseMatrix() Matrix4x4f[T] {
	rot := t.Rot

	rot.X = -rot.X
	rot.Y = -rot.Y
	rot.Z = -rot.Z

	x2 := rot.X * rot.X * 2
	y2 := rot.Y * rot.Y * 2
	z2 := rot.Z * rot.Z * 2

	xy := rot.X * rot.Y * 2
	wz := rot.W * rot.Z * 2
	xz := rot.X * rot.Z * 2
	wy := rot.W * rot.Y * 2
	yz := rot.Y * rot.Z * 2
	wx := rot.W * rot.X * 2

	m0 := Vector3f[T]{
		X: 1 - y2 - z2, Y: xy - wz, Z: xz + wy,
	}.ScaleInverseUniform(t.Scale.X)
	m1 := Vector3f[T]{
		X: xy + wz, Y: 1 - x2 - z2, Z: yz - wx,
	}.ScaleInverseUniform(t.Scale.Y)
	m2 := Vector3f[T]{
		X: xz - wy, Y: yz + wx, Z: 1 - x2 - y2,
	}.ScaleInverseUniform(t.Scale.Z)

	return Matrix4x4f[T]{
		{m0.X, m0.Y, m0.Z, t.Pivot.X - m0.Dot(Vector3f[T](t.Pos))},
		{m1.X, m1.Y, m1.Z, t.Pivot.Y - m1.Dot(Vector3f[T](t.Pos))},
		{m2.X, m2.Y, m2.Z, t.Pivot.Z - m2.Dot(Vector3f[T](t.Pos))},
		{0, 0, 0, 1},
	}
}
