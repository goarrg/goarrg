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
	"testing"
)

func TestC(t *testing.T) {
	c := PerspectiveCamera[float32]{
		SizeX: 3,
		SizeY: 3,
		FOV:   DegToRad[float32](90),
	}

	want := [][]Vector3f[float32]{
		{{-0.57735027, 0.57735027, 0.57735027}, {-0.70710678, 0, 0.70710678}, {-0.57735027, -0.57735027, 0.57735027}},
		{{0, 0.70710678, 0.70710678}, {0, 0, 1}, {0, -0.70710678, 0.70710678}},
		{{0.57735027, 0.57735027, 0.57735027}, {0.70710678, 0, 0.70710678}, {0.57735027, -0.57735027, 0.57735027}},
	}

	for x := 0; x < (int)(c.SizeX); x++ {
		for y := 0; y < (int)(c.SizeY); y++ {
			have := c.ScreenPointToRay(x, y).Dir
			if round(have) != round(want[x][y]) {
				t.Fatalf("[%d, %d] = %v != %v\n", x, y, have, want[x][y])
			}
		}
	}
}

func BenchmarkC(b *testing.B) {
	c := PerspectiveCamera[float32]{
		SizeX: 1920,
		SizeY: 1080,
		FOV:   90,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.ScreenPointToRay(1, 1)
	}
}
