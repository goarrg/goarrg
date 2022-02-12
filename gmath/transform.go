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
)

type Transform struct {
	Pos   Point3f64
	Rot   Quaternion
	Scale Vector3f64
}

func (t *Transform) TransformPoint(p Point3f64) Point3f64 {
	return t.Pos.Translate(t.Rot.Rotate(Vector3f64(p).Scale(t.Scale)))
}

func (t *Transform) TransformVector(v Vector3f64) Vector3f64 {
	return t.Rot.Rotate(v.Scale(t.Scale))
}

func (t *Transform) TransformDirection(v Vector3f64) Vector3f64 {
	return t.Rot.Rotate(v)
}

func (t *Transform) LookAt(up Vector3f64, target Point3f64) {
	m2 := t.Pos.VectorTo(target).Normalize()
	m0 := up.Cross(m2).Normalize()
	m1 := m2.Cross(m0)
	if trace := m0.X + m1.Y + m2.Z; trace > 0 {
		s := math.Sqrt(trace + 1)
		t.Rot.W = s * 0.5
		s = 0.5 / s
		t.Rot.X = (m1.Z - m2.Y) * s
		t.Rot.Y = (m2.X - m0.Z) * s
		t.Rot.Z = (m0.Y - m1.X) * s
		return
	}
	if (m0.X >= m1.Y) && (m0.X >= m2.Z) {
		s := math.Sqrt(1 + m0.X - m1.Y - m2.Z)
		t.Rot.X = 0.5 * s
		s = 0.5 / s
		t.Rot.Y = (m0.Y + m1.X) * s
		t.Rot.Z = (m0.Z + m2.X) * s
		t.Rot.W = (m1.Z - m2.Y) * s
		return
	}
	if m1.Y > m2.Z {
		s := math.Sqrt(1 + m1.Y - m0.X - m2.Z)
		t.Rot.Y = 0.5 * s
		s = 0.5 / s
		t.Rot.X = (m1.X + m0.Y) * s
		t.Rot.Z = (m2.Y + m1.Z) * s
		t.Rot.W = (m2.X - m0.Z) * s
		return
	}
	s := math.Sqrt(1 + m2.Z - m0.X - m1.Y)
	t.Rot.Z = 0.5 * s
	s = 0.5 / s
	t.Rot.X = (m2.X + m0.Z) * s
	t.Rot.Y = (m2.Y + m1.Z) * s
	t.Rot.W = (m0.Y - m1.X) * s
}

func (t *Transform) TranslationMatrix() Matrix4f64 {
	return Matrix4f64{
		{1, 0, 0, t.Pos.X},
		{0, 1, 0, t.Pos.Y},
		{0, 0, 1, t.Pos.Z},
		{0, 0, 0, 1},
	}
}

func (t *Transform) RotationMatrix() Matrix4f64 {
	x2 := t.Rot.X * t.Rot.X * 2
	y2 := t.Rot.Y * t.Rot.Y * 2
	z2 := t.Rot.Z * t.Rot.Z * 2

	xy := t.Rot.X * t.Rot.Y * 2
	wz := t.Rot.W * t.Rot.Z * 2
	xz := t.Rot.X * t.Rot.Z * 2
	wy := t.Rot.W * t.Rot.Y * 2
	yz := t.Rot.Y * t.Rot.Z * 2
	wx := t.Rot.W * t.Rot.X * 2

	return Matrix4f64{
		{1 - y2 - z2, xy - wz, xz + wy, 0},
		{xy + wz, 1 - x2 - z2, yz - wx, 0},
		{xz - wy, yz + wx, 1 - x2 - y2, 0},
		{0, 0, 0, 1},
	}
}

func (t *Transform) ScaleMatrix() Matrix4f64 {
	return Matrix4f64{
		{t.Scale.X, 0, 0, 0},
		{0, t.Scale.Y, 0, 0},
		{0, 0, t.Scale.Z, 0},
		{0, 0, 0, 1},
	}
}

func (t *Transform) ModelMatrix() Matrix4f64 {
	x2 := t.Rot.X * t.Rot.X * 2
	y2 := t.Rot.Y * t.Rot.Y * 2
	z2 := t.Rot.Z * t.Rot.Z * 2

	xy := t.Rot.X * t.Rot.Y * 2
	wz := t.Rot.W * t.Rot.Z * 2
	xz := t.Rot.X * t.Rot.Z * 2
	wy := t.Rot.W * t.Rot.Y * 2
	yz := t.Rot.Y * t.Rot.Z * 2
	wx := t.Rot.W * t.Rot.X * 2

	m0 := Vector3f64{
		X: 1 - y2 - z2, Y: xy - wz, Z: xz + wy,
	}.Scale(t.Scale)

	m1 := Vector3f64{
		X: xy + wz, Y: 1 - x2 - z2, Z: yz - wx,
	}.Scale(t.Scale)

	m2 := Vector3f64{
		X: xz - wy, Y: yz + wx, Z: 1 - x2 - y2,
	}.Scale(t.Scale)

	return Matrix4f64{
		{m0.X, m0.Y, m0.Z, t.Pos.X},
		{m1.X, m1.Y, m1.Z, t.Pos.Y},
		{m2.X, m2.Y, m2.Z, t.Pos.Z},
		{0, 0, 0, 1},
	}
}
