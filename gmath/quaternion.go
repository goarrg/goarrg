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

type Quaternion struct {
	X, Y, Z, W float64
}

// QuaternionFromEuler is creates a quaternion from eular rotation in the z x y order.
func QuaternionFromEuler(x, y, z float64) Quaternion {
	x *= 0.5
	y *= 0.5
	z *= 0.5

	cx := math.Cos(x)
	cy := math.Cos(y)
	cz := math.Cos(z)

	sx := math.Sin(x)
	sy := math.Sin(y)
	sz := math.Sin(z)

	cxcy := cx * cy
	cxsy := cx * sy
	sxsy := sx * sy
	sxcy := sx * cy

	return Quaternion{
		X: (sxcy * cz) + (cxsy * sz),
		Y: (cxsy * cz) - (sxcy * sz),
		Z: (cxcy * sz) - (sxsy * cz),
		W: (cxcy * cz) + (sxsy * sz),
	}
}

func QuaternionFromAngleAxis(t float64, axis Vector3f64) Quaternion {
	t *= 0.5
	st := math.Sin(t)

	return Quaternion{
		X: st * axis.X,
		Y: st * axis.Y,
		Z: st * axis.Z,
		W: math.Cos(t),
	}
}

func (q Quaternion) Multiply(q2 Quaternion) Quaternion {
	return Quaternion{
		X: (q.X * q2.W) + (q.Y * q2.Z) - (q.Z * q2.Y) + (q.W * q2.X),
		Y: (-q.X * q2.Z) + (q.Y * q2.W) + (q.Z * q2.X) + (q.W * q2.Y),
		Z: (q.X * q2.Y) - (q.Y * q2.X) + (q.Z * q2.W) + (q.W * q2.Z),
		W: (-q.X * q2.X) - (q.Y * q2.Y) - (q.Z * q2.Z) + (q.W * q2.W),
	}
}

func (q Quaternion) Rotate(v Vector3f64) Vector3f64 {
	x2 := q.X * q.X * 2
	y2 := q.Y * q.Y * 2
	z2 := q.Z * q.Z * 2

	xy := q.X * q.Y * 2
	wz := q.W * q.Z * 2
	xz := q.X * q.Z * 2
	wy := q.W * q.Y * 2
	yz := q.Y * q.Z * 2
	wx := q.W * q.X * 2

	m0 := Vector3f64{
		X: 1 - y2 - z2,
		Y: xy - wz,
		Z: xz + wy,
	}

	m1 := Vector3f64{
		X: xy + wz,
		Y: 1 - x2 - z2,
		Z: yz - wx,
	}

	m2 := Vector3f64{
		X: xz - wy,
		Y: yz + wx,
		Z: 1 - x2 - y2,
	}

	return Vector3f64{
		X: m0.Dot(v),
		Y: m1.Dot(v),
		Z: m2.Dot(v),
	}
}

func (q Quaternion) Magnitude() float64 {
	return math.Sqrt((q.X * q.X) + (q.Y * q.Y) + (q.Z * q.Z) + (q.W * q.W))
}

func (q Quaternion) Normalize() Quaternion {
	m := q.Magnitude()

	if m == 0 || m == 1 {
		return q
	}

	return Quaternion{
		X: q.X / m,
		Y: q.Y / m,
		Z: q.Z / m,
		W: q.W / m,
	}
}
