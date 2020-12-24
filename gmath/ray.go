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

type Ray struct {
	Dir    Vector3f64
	InvDir Vector3f64
	Org    Point3f64
}

func (r *Ray) IntersectionAABB(p Point3f64, b Bounds3f64) float64 {
	t1 := (p.X - b.Min.X - r.Org.X) * r.InvDir.X
	t2 := (p.X + b.Max.X - r.Org.X) * r.InvDir.X

	tmin := math.Min(t1, t2)
	tmax := math.Max(t1, t2)

	t1 = (p.Y - b.Min.Y - r.Org.Y) * r.InvDir.Y
	t2 = (p.Y + b.Max.Y - r.Org.Y) * r.InvDir.Y

	tmin = math.Max(tmin, math.Min(math.Min(t1, t2), tmax))
	tmax = math.Min(tmax, math.Max(math.Max(t1, t2), tmin))

	t1 = (p.Z - b.Min.Z - r.Org.Z) * r.InvDir.Z
	t2 = (p.Z + b.Max.Z - r.Org.Z) * r.InvDir.Z

	tmin = math.Max(math.Max(tmin, math.Min(math.Min(t1, t2), tmax)), 0.0)
	tmax = math.Min(tmax, math.Max(math.Max(t1, t2), tmin))

	if tmax > tmin {
		return tmin
	}

	return math.NaN()
}
