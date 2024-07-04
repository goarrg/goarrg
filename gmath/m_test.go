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

func (m Matrix4x4f[T]) MultiplyPoint(p Point3f[T]) Point3f[T] {
	return Point3f[T]{
		X: m[0][0]*p.X + m[0][1]*p.Y + m[0][2]*p.Z + m[0][3],
		Y: m[1][0]*p.X + m[1][1]*p.Y + m[1][2]*p.Z + m[1][3],
		Z: m[2][0]*p.X + m[2][1]*p.Y + m[2][2]*p.Z + m[2][3],
	}
}

func TestMatrix(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform[float32]{}
		target := Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}

		tr.Pivot = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
		tr.Pos = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
		tr.Rot = QuaternionFromEuler(rand.Float32(), rand.Float32(), rand.Float32())
		tr.Scale = Vector3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}

		want := Vector3f[float32](tr.TransformPoint(target))
		got := Vector3f[float32](tr.ModelMatrix().MultiplyPoint(target))

		if want.Subtract(got).Magnitude() > 1e-6 {
			t.Fatalf("Want\n%+v\nbut got\n%+v", want, got)
		}
	}
}

func TestMatrixTranspose(t *testing.T) {
	tr := Transform[float32]{Pos: Point3f[float32]{X: 2, Y: 3, Z: 4}}
	got := tr.TranslationMatrix().Transpose()
	want := Matrix4x4f[float32]{
		[4]float32{1, 0, 0, 0},
		[4]float32{0, 1, 0, 0},
		[4]float32{0, 0, 1, 0},
		[4]float32{2, 3, 4, 1},
	}

	if got != want {
		t.Fatalf("Want\n%+v\nbut got\n%+v", want, got)
	}
}

func TestMatrixInvert(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform[float32]{}
		tr.Pivot = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
		tr.Pos = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
		tr.Rot = QuaternionFromEuler(rand.Float32(), rand.Float32(), rand.Float32())
		tr.Scale = Vector3f[float32]{1, 1, 1}

		m := tr.ModelMatrix()
		mI := m.Invert()

		got := m.Multiply(mI)
		want := Matrix4x4f[float32]{
			[4]float32{1, 0, 0, 0},
			[4]float32{0, 1, 0, 0},
			[4]float32{0, 0, 1, 0},
			[4]float32{0, 0, 0, 1},
		}

		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if math.Abs(float64(got[i][j]-want[i][j])) > 1e-6 {
					t.Fatalf("Want\n%+v\nbut got\n%+v", want, got)
				}
			}
		}
	}
}

func TestMatrixModel(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform[float32]{}

		tr.Pivot = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
		tr.Pos = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
		tr.Rot = QuaternionFromEuler(rand.Float32(), rand.Float32(), rand.Float32())
		tr.Scale = Vector3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}

		m := tr.ModelMatrix()
		// test ModelMatrix()
		{
			want := tr.TranslationMatrix().Multiply(tr.RotationMatrix().Multiply(tr.ScaleMatrix().Multiply(tr.PivotMatrix())))

			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if math.Abs(float64(want[i][j]-m[i][j])) > 1e-6 {
						t.Fatalf("Want\n%+v\nbut got\n%+v", want, m)
					}
				}
			}
		}

		// test ModelInverseMatrix()
		{
			m2 := tr.ModelInverseMatrix()
			got := m.Multiply(m2)
			want := Matrix4x4f[float32]{
				[4]float32{1, 0, 0, 0},
				[4]float32{0, 1, 0, 0},
				[4]float32{0, 0, 1, 0},
				[4]float32{0, 0, 0, 1},
			}

			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if math.Abs(float64(got[i][j]-want[i][j])) > 1e-6 {
						t.Fatalf("Want\n%+v\nbut got\n%+v", want, got)
					}
				}
			}
		}
	}
}

func TestMatrixView(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		tr := Transform[float32]{}

		tr.Pos = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
		tr.Rot = QuaternionFromEuler(rand.Float32(), rand.Float32(), rand.Float32())

		c := PerspectiveCamera[float32]{Transform: tr}

		m := c.ViewMatrix()
		// test ViewMatrix()
		{
			tr := tr
			tr.Rot.X = -tr.Rot.X
			tr.Rot.Y = -tr.Rot.Y
			tr.Rot.Z = -tr.Rot.Z

			rm := tr.RotationMatrix()
			tm := tr.TranslationMatrix()
			want := rm.Multiply(tm)

			want[0][3] = -want[0][3]
			want[1][3] = -want[1][3]
			want[2][3] = -want[2][3]

			if want != m {
				t.Fatalf("Want\n%+v\nbut got\n%+v", want, m)
			}
		}

		// test ViewInverseMatrix()
		{
			m2 := c.ViewInverseMatrix()
			got := m.Multiply(m2)
			want := Matrix4x4f[float32]{
				[4]float32{1, 0, 0, 0},
				[4]float32{0, 1, 0, 0},
				[4]float32{0, 0, 1, 0},
				[4]float32{0, 0, 0, 1},
			}

			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if math.Abs(float64(got[i][j]-want[i][j])) > 1e-6 {
						t.Fatalf("Want\n%+v\nbut got\n%+v", want, got)
					}
				}
			}
		}
	}
}

func TestMatrixPerspective(t *testing.T) {
	for i := 0; i < 1e6; i++ {
		c := PerspectiveCamera[float32]{
			SizeX: rand.Float32() * 1920, SizeY: rand.Float32() * 1080,
			FOV: math.Pi * rand.Float32(), ZNear: rand.Float32(),
		}
		m := c.ProjectionMatrix()
		m2 := c.ProjectionInverseMatrix()
		got := m.Multiply(m2)
		want := Matrix4x4f[float32]{
			[4]float32{1, 0, 0, 0},
			[4]float32{0, 1, 0, 0},
			[4]float32{0, 0, 1, 0},
			[4]float32{0, 0, 0, 1},
		}

		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if math.Abs(float64(got[i][j]-want[i][j])) > 1e-6 {
					t.Fatalf("Want\n%+v\nbut got\n%+v", want, got)
				}
			}
		}
	}
}

var (
	globalP Point3f32
	globalM Matrix4x4f32
)

func BenchmarkMatrix(b *testing.B) {
	tr := Transform[float32]{}
	target := Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}

	tr.Pos = Point3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}
	tr.Rot = QuaternionFromEuler(rand.Float32(), rand.Float32(), rand.Float32())
	tr.Scale = Vector3f[float32]{rand.Float32(), rand.Float32(), rand.Float32()}

	c := PerspectiveCamera[float32]{Transform: tr, FOV: rand.Float32(), SizeX: rand.Float32(), SizeY: rand.Float32()}

	var p Point3f32
	var m Matrix4x4f32

	b.Run("TransformPoint", func(b *testing.B) {
		b.Run("Slow", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p = tr.ModelMatrix().MultiplyPoint(target)
			}
		})
		globalP = p
		b.Run("Optimized", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p = tr.TransformPoint(target)
			}
		})
		globalP = p
	})

	b.Run("Model", func(b *testing.B) {
		b.Run("Slow", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m = tr.TranslationMatrix().Multiply(tr.RotationMatrix().Multiply(tr.ScaleMatrix()))
			}
		})
		globalM = m
		b.Run("Optimized", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m = tr.ModelMatrix()
			}
		})
		globalM = m
		b.Run("Inverse", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m = tr.ModelInverseMatrix()
			}
		})
		globalM = m
	})

	b.Run("View", func(b *testing.B) {
		b.Run("Slow", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tr.Rot.X = -tr.Rot.X
				tr.Rot.Y = -tr.Rot.Y
				tr.Rot.Z = -tr.Rot.Z

				rm := tr.RotationMatrix()
				tm := tr.TranslationMatrix()
				m = rm.Multiply(tm)

				m[0][3] = -m[0][3]
				m[1][3] = -m[1][3]
				m[2][3] = -m[2][3]
			}
		})
		globalM = m
		b.Run("Optimized", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m = c.ViewMatrix()
			}
		})
		globalM = m
		b.Run("Inverse", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m = c.ViewInverseMatrix()
			}
		})
		globalM = m
	})

	m2 := tr.ModelMatrix()

	b.Run("MultiplyPoint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p = m2.MultiplyPoint(target)
		}
	})
	globalP = p

	b.Run("ModelMultiplyPoint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p = tr.ModelMatrix().MultiplyPoint(target)
		}
	})
	globalP = p

	b.Run("Multiply", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m = m2.Multiply(m)
		}
	})
	globalM = m

	b.Run("Transpose", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m = m2.Transpose()
		}
	})
	globalM = m

	b.Run("Invert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m = m2.Invert()
		}
	})
	globalM = m

	b.Run("ToArrayf32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m = m2.ToArrayf32()
		}
	})
	globalM = m

	b.Run("TransposedToArrayf32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m = m2.TransposedToArrayf32()
		}
	})
	globalM = m
}
