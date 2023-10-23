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

type Camera struct {
	Transform

	SizeX float64
	SizeY float64
	FOV   float64
}

func (c *Camera) ScreenPointToRay(x, y int) Ray {
	dir := Vector3f64{
		X: float64(x) - ((c.SizeX - 1) / 2),
		Y: ((c.SizeY - 1) / 2) - float64(y),
		Z: ((c.SizeY - 1) / 2) / math.Tan(c.FOV*0.5),
	}.Normalize()

	dir = c.TransformDirection(dir)
	invDir := Vector3f64{
		X: 1 / dir.X,
		Y: 1 / dir.Y,
		Z: 1 / dir.Z,
	}

	return Ray{
		Dir:    dir,
		InvDir: invDir,
		Org:    c.Pos,
	}
}

func (c *Camera) ViewMatrix() Matrix4f64 {
	rot := c.Transform.Rot

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

	m0 := Vector3f64{
		X: 1 - y2 - z2, Y: xy - wz, Z: xz + wy,
	}

	m1 := Vector3f64{
		X: xy + wz, Y: 1 - x2 - z2, Z: yz - wx,
	}

	m2 := Vector3f64{
		X: xz - wy, Y: yz + wx, Z: 1 - x2 - y2,
	}

	return Matrix4f64{
		{m0.X, m0.Y, m0.Z, -(m0.Dot(Vector3f64(c.Transform.Pos)))},
		{m1.X, m1.Y, m1.Z, -(m1.Dot(Vector3f64(c.Transform.Pos)))},
		{m2.X, m2.Y, m2.Z, -(m2.Dot(Vector3f64(c.Transform.Pos)))},
		{0, 0, 0, 1},
	}
}

func (c *Camera) PerspectiveMatrix() Matrix4f64 {
	t := math.Tan(c.FOV * 0.5)

	return Matrix4f64{
		{1 / ((c.SizeX / c.SizeY) * t), 0, 0, 0},
		{0, 1 / t, 0, 0},
		{0, 0, 0, 1},
		{0, 0, 1, 0},
	}
}

func (c *Camera) PerspectiveInverseMatrix() Matrix4f64 {
	t := math.Tan(c.FOV * 0.5)

	return Matrix4f64{
		{((c.SizeX / c.SizeY) * t), 0, 0, 0},
		{0, t, 0, 0},
		{0, 0, 0, 1},
		{0, 0, 1, 0},
	}
}
