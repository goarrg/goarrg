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
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkQ(b *testing.B) {
	b.Run("Euler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			QuaternionFromEuler[float32](1, 2, 3)
		}
	})

	b.Run("Axis", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			QuaternionFromAngleAxis(1, Vector3f[float32]{2, 3, 4})
		}
	})

	b.Run("Rotate", func(b *testing.B) {
		q := QuaternionFromAngleAxis(1, Vector3f[float32]{2, 3, 4})
		v := Vector3f[float32]{5, 6, 7}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = q.Rotate(v)
		}
	})

	b.Run("RotateP", func(b *testing.B) {
		q := QuaternionFromAngleAxis(1, Vector3f[float32]{2, 3, 4})
		v := Point3f[float32]{5, 6, 7}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = q.Rotate(Vector3f[float32](v))
		}
	})

	q1 := QuaternionFromAngleAxis(1, Vector3f[float32]{2, 3, 4})

	q2 := QuaternionFromAngleAxis(5, Vector3f[float32]{6, 7, 8})

	b.Run("Multiply", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = q1.Multiply(q2)
		}
	})
}

func TestQ(t *testing.T) {
	failed := int32(0)
	count := runtime.NumCPU() - 2
	chunkSize := float32(math.Ceil(361 / (float64)(count)))
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			for x := float32(0.0); x < chunkSize; x++ {
				x := x + (float32(i) * chunkSize) - 180
				if x > 180 {
					return
				}
				for y := float32(-180.0); y <= 180; y++ {
					for z := float32(-180.0); z <= 180; z++ {
						if atomic.LoadInt32(&failed) != 0 {
							return
						}
						qEuler := QuaternionFromEuler(x, y, z)

						qx := QuaternionFromAngleAxis(x, Vector3f[float32]{1, 0, 0})
						qy := QuaternionFromAngleAxis(y, Vector3f[float32]{0, 1, 0})
						qz := QuaternionFromAngleAxis(z, Vector3f[float32]{0, 0, 1})

						qzxy := qy.Multiply(qx).Multiply(qz)

						if qzxy != qEuler {
							atomic.StoreInt32(&failed, 1)
							t.Logf("%f %f %f, %+v != %+v\n", x, y, z, qEuler, qzxy)
						}
					}
				}
			}
		}()
	}
	wg.Wait()

	if failed != 0 {
		t.FailNow()
	}
}

func q(x, y, z float32) Quaternion[float32] {
	q := QuaternionFromEuler(x*(math.Pi/180), y*(math.Pi/180), z*(math.Pi/180))
	return q
}

func round(v Vector3f[float32]) Vector3f[float32] {
	return Vector3f[float32]{
		X: float32(math.Round(float64(v.X*1e6)) / 1e6),
		Y: float32(math.Round(float64(v.Y*1e6)) / 1e6),
		Z: float32(math.Round(float64(v.Z*1e6)) / 1e6),
	}
}

func TestQRot(t *testing.T) {
	q := []Quaternion[float32]{
		q(90, 0, 0),
		q(-90, 0, 0),
		q(0, 90, 0),
		q(0, -90, 0),
		q(0, 0, 90),
		q(0, 0, -90),
		q(70, 80, 90),
		q(70, 80, 90),
	}

	v := []Vector3f[float32]{
		{0, 0, 1},
		{0, 0, 1},
		{0, 0, 1},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 0},
		{0, 0, 1},
		{0, 1, 0},
	}

	want := []Vector3f[float32]{
		{0, -1, 0},
		{0, 1, 0},
		{1, 0, 0},
		{-1, 0, 0},
		{-1, 0, 0},
		{1, 0, 0},
		{0.33682409, -0.93969262, 0.05939117},
		{-0.17364818, 0, 0.98480775},
	}

	if len(q) != len(v) && len(v) != len(want) {
		t.Fatal("q, v and want must be of same length")
	}

	for i := 0; i < len(q); i++ {
		if have := q[i].Rotate(v[i]); round(have) != round(want[i]) {
			t.Fatalf("[%d] %+v * %+v = %+v not %+v", i, q[i], v[i], have, want[i])
		}
	}
}
