//+build !debug

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
	"math/rand"
	"testing"
)

func TestLookAt(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		p1 := Point3f64{
			X: rand.Float64(),
			Y: rand.Float64(),
			Z: rand.Float64(),
		}

		p2 := Point3f64{
			X: rand.Float64(),
			Y: rand.Float64(),
			Z: rand.Float64(),
		}

		tr := Transform{Pos: p1}
		tr.LookAt(Vector3f64{Y: 1}, p2)

		f := tr.TransformDirection(Vector3f64{Z: 1})
		d := p1.VectorTo(p2).Normalize()

		if f.Subtract(d).Magnitude() > 0.000001 {
			t.Logf("Want: %f but got %f\n", d, f)
			t.FailNow()
		}
	}
}

func BenchmarkLookAt(b *testing.B) {
	t := Transform{}

	for i := 0; i < b.N; i++ {
		t.LookAt(Vector3f64{Y: 1}, Point3f64{})
	}
}
