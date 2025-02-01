//go:build !goarrg_build_debug
// +build !goarrg_build_debug

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
		p1 := Point3f[float32]{
			X: rand.Float32(),
			Y: rand.Float32(),
			Z: rand.Float32(),
		}

		p2 := Point3f[float32]{
			X: rand.Float32(),
			Y: rand.Float32(),
			Z: rand.Float32(),
		}

		tr := Transform[float32]{Pos: p1}
		tr.LookAt(Vector3f[float32]{Y: 1}, p2)

		f := tr.TransformDirection(Vector3f[float32]{Z: 1})
		d := p1.VectorTo(p2).Normalize()

		if f.Subtract(d).Magnitude() > 0.000001 {
			t.Logf("Want: %f but got %f\n", d, f)
			t.FailNow()
		}
	}
}

func BenchmarkLookAt(b *testing.B) {
	t := Transform[float32]{}

	for i := 0; i < b.N; i++ {
		t.LookAt(Vector3f[float32]{Y: 1}, Point3f[float32]{})
	}
}
