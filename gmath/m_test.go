//go:build !debug
// +build !debug

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
	"math/rand"
	"testing"
)

func (m Matrix4f64) MultiplyPoint(p Point3f64) Point3f64 {
	return Point3f64{
		X: m[0][0]*p.X + m[0][1]*p.Y + m[0][2]*p.Z + m[0][3],
		Y: m[1][0]*p.X + m[1][1]*p.Y + m[1][2]*p.Z + m[1][3],
		Z: m[2][0]*p.X + m[2][1]*p.Y + m[2][2]*p.Z + m[2][3],
	}
}

func TestMatrix(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform{}
		target := Point3f64{rand.Float64(), rand.Float64(), rand.Float64()}

		tr.Pos = Point3f64{rand.Float64(), rand.Float64(), rand.Float64()}
		tr.Rot = QuaternionFromEuler(rand.Float64(), rand.Float64(), rand.Float64())
		tr.Scale = Vector3f64{rand.Float64(), rand.Float64(), rand.Float64()}

		tp := Vector3f64(tr.TransformPoint(target))
		mm := Vector3f64(tr.ModelMatrix().MultiplyPoint(target))

		if tp.Subtract(mm).Magnitude() > 0.00000001 {
			t.Fatal(tp, mm)
		}
	}
}

func TestMatrixTranspose(t *testing.T) {
	tr := Transform{Pos: Point3f64{X: 2, Y: 3, Z: 4}}
	got := tr.TranslationMatrix().Transpose()
	want := Matrix4f64{
		[4]float64{1, 0, 0, 0},
		[4]float64{0, 1, 0, 0},
		[4]float64{0, 0, 1, 0},
		[4]float64{2, 3, 4, 1},
	}

	if got != want {
		t.Fatalf("Want %+v but got %+v", want, got)
	}
}

func TestMatrixInvert(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform{}
		tr.Pos = Point3f64{rand.Float64(), rand.Float64(), rand.Float64()}
		tr.Rot = QuaternionFromEuler(rand.Float64(), rand.Float64(), rand.Float64()).Normalize()
		tr.Scale = Vector3f64{1, 1, 1}

		m := tr.ModelMatrix()
		mI := m.Invert()

		got := m.Multiply(mI)
		want := Matrix4f64{
			[4]float64{1, 0, 0, 0},
			[4]float64{0, 1, 0, 0},
			[4]float64{0, 0, 1, 0},
			[4]float64{0, 0, 0, 1},
		}

		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if math.Abs(got[i][j]-want[i][j]) > 1e-15 {
					t.Fatalf("Want %+v but got %+v", want, got)
				}
			}
		}
	}
}

func TestMatrixModel(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform{}

		tr.Pos = Point3f64{rand.Float64(), rand.Float64(), rand.Float64()}
		tr.Rot = QuaternionFromEuler(rand.Float64(), rand.Float64(), rand.Float64())
		tr.Scale = Vector3f64{rand.Float64(), rand.Float64(), rand.Float64()}

		m := tr.ModelMatrix()
		m2 := tr.TranslationMatrix().Multiply(tr.RotationMatrix().Multiply(tr.ScaleMatrix()))

		if m != m2 {
			t.Fatal(m, m2)
		}
	}
}

func TestMatrixView(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform{}

		tr.Pos = Point3f64{rand.Float64(), rand.Float64(), rand.Float64()}
		tr.Rot = QuaternionFromEuler(rand.Float64(), rand.Float64(), rand.Float64())

		c := Camera{Transform: tr}

		tr.Rot.X = -tr.Rot.X
		tr.Rot.Y = -tr.Rot.Y
		tr.Rot.Z = -tr.Rot.Z

		rm := tr.RotationMatrix()
		tm := tr.TranslationMatrix()
		m := rm.Multiply(tm)

		m[0][3] = -m[0][3]
		m[1][3] = -m[1][3]
		m[2][3] = -m[2][3]

		m2 := c.ViewMatrix()

		if m != m2 {
			t.Fatal("\n", m, "\n", m2)
		}
	}
}

func BenchmarkMatrix(b *testing.B) {
	tr := Transform{}
	target := Point3f64{rand.Float64(), rand.Float64(), rand.Float64()}

	tr.Pos = Point3f64{rand.Float64(), rand.Float64(), rand.Float64()}
	tr.Rot = QuaternionFromEuler(rand.Float64(), rand.Float64(), rand.Float64())
	tr.Scale = Vector3f64{rand.Float64(), rand.Float64(), rand.Float64()}

	c := Camera{Transform: tr, FOV: rand.Float64(), SizeX: rand.Float64(), SizeY: rand.Float64()}

	b.Run("TransformPoint", func(b *testing.B) {
		b.Run("Slow", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tr.ModelMatrix().MultiplyPoint(target)
			}
		})
		b.Run("Optimized", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tr.TransformPoint(target)
			}
		})
	})

	b.Run("Model", func(b *testing.B) {
		b.Run("Slow", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tr.TranslationMatrix().Multiply(tr.RotationMatrix().Multiply(tr.ScaleMatrix()))
			}
		})
		b.Run("Optimized", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tr.ModelMatrix()
			}
		})
	})

	b.Run("View", func(b *testing.B) {
		b.Run("Slow", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tr.Rot.X = -tr.Rot.X
				tr.Rot.Y = -tr.Rot.Y
				tr.Rot.Z = -tr.Rot.Z

				rm := tr.RotationMatrix()
				tm := tr.TranslationMatrix()
				m := rm.Multiply(tm)

				m[0][3] = -m[0][3]
				m[1][3] = -m[1][3]
				m[2][3] = -m[2][3]
			}
		})
		b.Run("Optimized", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = c.ViewMatrix()
			}
		})
	})

	m := tr.ModelMatrix()

	b.Run("MultiplyPoint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m.MultiplyPoint(target)
		}
	})

	b.Run("ModelMultiplyPoint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = tr.ModelMatrix().MultiplyPoint(target)
		}
	})

	b.Run("Multiply", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m.Multiply(m)
		}
	})

	b.Run("Transpose", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m.Transpose()
		}
	})

	b.Run("Invert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m.Invert()
		}
	})
}
