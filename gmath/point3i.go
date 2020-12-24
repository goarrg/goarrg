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

type Point3i struct {
	X, Y, Z int
}

func Point3iFromArray(a [3]int) Point3i {
	return Point3i{
		X: a[0],
		Y: a[1],
		Z: a[2],
	}
}

func (p Point3i) VectorTo(p2 Point3i) Vector3i {
	return Vector3i{
		X: p2.X - p.X,
		Y: p2.Y - p.Y,
		Z: p2.Z - p.Z,
	}
}

func (p Point3i) Translate(v Vector3i) Point3i {
	return Point3i{
		X: p.X + v.X,
		Y: p.Y + v.Y,
		Z: p.Z + v.Z,
	}
}

func (p Point3i) ToArray() [3]int {
	return [3]int{p.X, p.Y, p.Z}
}
