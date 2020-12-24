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

type Vector3f64 struct {
	X, Y, Z float64
}

func Vector3f64FromArray(a [3]float64) Vector3f64 {
	return Vector3f64{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (v Vector3f64) IsNAN() bool {
	return math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsNaN(v.Z)
}

func (v Vector3f64) Add(v2 Vector3f64) Vector3f64 {
	return Vector3f64{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vector3f64) AddUniform(f float64) Vector3f64 {
	return Vector3f64{
		X: v.X + f,
		Y: v.Y + f,
		Z: v.Z + f,
	}
}

func (v Vector3f64) Subtract(v2 Vector3f64) Vector3f64 {
	return Vector3f64{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

func (v Vector3f64) SubtractUniform(f float64) Vector3f64 {
	return Vector3f64{
		X: v.X - f,
		Y: v.Y - f,
		Z: v.Z - f,
	}
}

func (v Vector3f64) Scale(v2 Vector3f64) Vector3f64 {
	return Vector3f64{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
		Z: v.Z * v2.Z,
	}
}

func (v Vector3f64) ScaleInverse(v2 Vector3f64) Vector3f64 {
	return Vector3f64{
		X: v.X / v2.X,
		Y: v.Y / v2.Y,
		Z: v.Z / v2.Z,
	}
}

func (v Vector3f64) ScaleUniform(s float64) Vector3f64 {
	return Vector3f64{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func (v Vector3f64) ScaleInverseUniform(s float64) Vector3f64 {
	return Vector3f64{
		X: v.X / s,
		Y: v.Y / s,
		Z: v.Z / s,
	}
}

func (v Vector3f64) Dot(v2 Vector3f64) float64 {
	return (v.X * v2.X) + (v.Y * v2.Y) + (v.Z * v2.Z)
}

func (v Vector3f64) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z))
}

func (v Vector3f64) Angle(v2 Vector3f64) float64 {
	dot := v.Dot(v2)
	ma := v.Magnitude()
	mb := v2.Magnitude()

	return math.Acos(dot / (ma * mb))
}

func (v Vector3f64) AngleAxis(v2 Vector3f64) (float64, Vector3f64) {
	return v.Angle(v2), v.Normalize().Cross(v2.Normalize()).Normalize()
}

func (v Vector3f64) Normalize() Vector3f64 {
	m := v.Magnitude()

	if m == 0 || m == 1 {
		return v
	}

	return Vector3f64{
		X: v.X / m,
		Y: v.Y / m,
		Z: v.Z / m,
	}
}

func (v Vector3f64) Cross(v2 Vector3f64) Vector3f64 {
	return Vector3f64{
		X: (v.Y * v2.Z) - (v.Z * v2.Y),
		Y: (v.Z * v2.X) - (v.X * v2.Z),
		Z: (v.X * v2.Y) - (v.Y * v2.X),
	}
}

func (v Vector3f64) OrthoNormalize(tangent Vector3f64) Vector3f64 {
	v = v.Normalize()
	tangent = tangent.Normalize()

	return tangent.Subtract(v.ScaleUniform(tangent.Dot(v))).Normalize()
}

func (v Vector3f64) ToArray() [3]float64 {
	return [3]float64{v.X, v.Y, v.Z}
}

func (v Vector3f64) ToArrayf32() [3]float32 {
	return [3]float32{float32(v.X), float32(v.Y), float32(v.Z)}
}
