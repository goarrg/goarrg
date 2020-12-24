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
	"testing"
)

func BenchmarkR(b *testing.B) {
	c := Camera{
		SizeX: 1920,
		SizeY: 1080,
		FOV:   90,
	}

	r := c.ScreenPointToRay(1, 1)
	aabb := Bounds3f64{
		Min: Vector3f64{-100, -100, -100},
		Max: Vector3f64{100, 100, 100},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = r.IntersectionAABB(c.Pos, aabb)
	}
}
