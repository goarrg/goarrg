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

type Camera[T constraints.Float] interface {
	ScreenPointToRay(x, y int) Ray[T]

	ViewMatrix() Matrix4x4f[T]
	ViewInverseMatrix() Matrix4x4f[T]

	ProjectionMatrix() Matrix4x4f[T]
	ProjectionInverseMatrix() Matrix4x4f[T]
}

type (
	Cameraf32 = Camera[float32]
	Cameraf64 = Camera[float64]
)

func viewMatrix[T constraints.Float](transform Transform[T]) Matrix4x4f[T] {
	rot := transform.Rot

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
	}
	m1 := Vector3f[T]{
		X: xy + wz, Y: 1 - x2 - z2, Z: yz - wx,
	}
	m2 := Vector3f[T]{
		X: xz - wy, Y: yz + wx, Z: 1 - x2 - y2,
	}

	return Matrix4x4f[T]{
		{m0.X, m0.Y, m0.Z, -(m0.Dot(Vector3f[T](transform.Pos)))},
		{m1.X, m1.Y, m1.Z, -(m1.Dot(Vector3f[T](transform.Pos)))},
		{m2.X, m2.Y, m2.Z, -(m2.Dot(Vector3f[T](transform.Pos)))},
		{0, 0, 0, 1},
	}
}

func viewInverseMatrix[T constraints.Float](transform Transform[T]) Matrix4x4f[T] {
	x2 := transform.Rot.X * transform.Rot.X * 2
	y2 := transform.Rot.Y * transform.Rot.Y * 2
	z2 := transform.Rot.Z * transform.Rot.Z * 2

	xy := transform.Rot.X * transform.Rot.Y * 2
	wz := transform.Rot.W * transform.Rot.Z * 2
	xz := transform.Rot.X * transform.Rot.Z * 2
	wy := transform.Rot.W * transform.Rot.Y * 2
	yz := transform.Rot.Y * transform.Rot.Z * 2
	wx := transform.Rot.W * transform.Rot.X * 2

	m0 := Vector3f[T]{
		X: 1 - y2 - z2, Y: xy - wz, Z: xz + wy,
	}
	m1 := Vector3f[T]{
		X: xy + wz, Y: 1 - x2 - z2, Z: yz - wx,
	}
	m2 := Vector3f[T]{
		X: xz - wy, Y: yz + wx, Z: 1 - x2 - y2,
	}

	return Matrix4x4f[T]{
		{m0.X, m0.Y, m0.Z, transform.Pos.X},
		{m1.X, m1.Y, m1.Z, transform.Pos.Y},
		{m2.X, m2.Y, m2.Z, transform.Pos.Z},
		{0, 0, 0, 1},
	}
}

type PerspectiveCamera[T constraints.Float] struct {
	Transform[T]

	SizeX T
	SizeY T
	FOV   T

	ZNear T
}

type (
	PerspectiveCameraf32 = PerspectiveCamera[float32]
	PerspectiveCameraf64 = PerspectiveCamera[float64]
)

var (
	_ Camera[float32] = &PerspectiveCameraf32{}
	_ Camera[float64] = &PerspectiveCameraf64{}
)

func (c *PerspectiveCamera[T]) ScreenPointToRay(x, y int) Ray[T] {
	dir := Vector3f[T]{
		X: T(x) - ((c.SizeX - 1) / 2),
		Y: ((c.SizeY - 1) / 2) - T(y),
		Z: ((c.SizeY - 1) / 2) / T(math.Tan(float64(c.FOV*0.5))),
	}.Normalize()

	dir = c.TransformDirection(dir)
	invDir := Vector3f[T]{
		X: 1 / dir.X,
		Y: 1 / dir.Y,
		Z: 1 / dir.Z,
	}

	return Ray[T]{
		Dir:    dir,
		InvDir: invDir,
		Org:    c.Pos,
	}
}

func (c *PerspectiveCamera[T]) ViewMatrix() Matrix4x4f[T] {
	return viewMatrix(c.Transform)
}

func (c *PerspectiveCamera[T]) ViewInverseMatrix() Matrix4x4f[T] {
	return viewInverseMatrix(c.Transform)
}

func (c *PerspectiveCamera[T]) ProjectionMatrix() Matrix4x4f[T] {
	t := T(math.Tan(float64(c.FOV * 0.5)))

	return Matrix4x4f[T]{
		{1 / ((c.SizeX / c.SizeY) * t), 0, 0, 0},
		{0, 1 / t, 0, 0},
		{0, 0, 0, c.ZNear},
		{0, 0, 1, 0},
	}
}

func (c *PerspectiveCamera[T]) ProjectionInverseMatrix() Matrix4x4f[T] {
	t := T(math.Tan(float64(c.FOV * 0.5)))

	return Matrix4x4f[T]{
		{((c.SizeX / c.SizeY) * t), 0, 0, 0},
		{0, t, 0, 0},
		{0, 0, 0, 1},
		{0, 0, 1 / c.ZNear, 0},
	}
}
