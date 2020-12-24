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

type Vector3i struct {
	X, Y, Z int
}

func Vector3iFromArray(a [3]int) Vector3i {
	return Vector3i{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (v Vector3i) Add(v2 Vector3i) Vector3i {
	return Vector3i{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vector3i) AddUniform(i int) Vector3i {
	return Vector3i{
		X: v.X + i,
		Y: v.Y + i,
		Z: v.Z + i,
	}
}

func (v Vector3i) Subtract(v2 Vector3i) Vector3i {
	return Vector3i{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

func (v Vector3i) SubtractUniform(i int) Vector3i {
	return Vector3i{
		X: v.X - i,
		Y: v.Y - i,
		Z: v.Z - i,
	}
}

func (v Vector3i) Scale(v2 Vector3i) Vector3i {
	return Vector3i{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
		Z: v.Z * v2.Z,
	}
}

func (v Vector3i) ScaleInverse(v2 Vector3i) Vector3i {
	return Vector3i{
		X: v.X / v2.X,
		Y: v.Y / v2.Y,
		Z: v.Z / v2.Z,
	}
}

func (v Vector3i) ScaleUniform(s int) Vector3i {
	return Vector3i{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func (v Vector3i) ScaleInverseUniform(s int) Vector3i {
	return Vector3i{
		X: v.X / s,
		Y: v.Y / s,
		Z: v.Z / s,
	}
}

func (v Vector3i) ToArray() [3]int {
	return [3]int{v.X, v.Y, v.Z}
}
