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

import "golang.org/x/exp/constraints"

type Matrix4x4f[T constraints.Float] [4][4]T

type (
	Matrix4x4f32 = Matrix4x4f[float32]
	Matrix4x4f64 = Matrix4x4f[float64]
)

func (m Matrix4x4f[T]) Multiply(m2 Matrix4x4f[T]) Matrix4x4f[T] {
	mOut := Matrix4x4f[T]{}

	for i := 0; i < 4; i++ {
		mOut[i] = [4]T{
			m[i][0]*m2[0][0] + m[i][1]*m2[1][0] + m[i][2]*m2[2][0] + m[i][3]*m2[3][0],
			m[i][0]*m2[0][1] + m[i][1]*m2[1][1] + m[i][2]*m2[2][1] + m[i][3]*m2[3][1],
			m[i][0]*m2[0][2] + m[i][1]*m2[1][2] + m[i][2]*m2[2][2] + m[i][3]*m2[3][2],
			m[i][0]*m2[0][3] + m[i][1]*m2[1][3] + m[i][2]*m2[2][3] + m[i][3]*m2[3][3],
		}
	}

	return mOut
}

func (m Matrix4x4f[T]) Transpose() Matrix4x4f[T] {
	mOut := Matrix4x4f[T]{}

	for i := 0; i < 4; i++ {
		mOut[i] = [4]T{
			m[0][i],
			m[1][i],
			m[2][i],
			m[3][i],
		}
	}

	return mOut
}

func (m Matrix4x4f[T]) Invert() Matrix4x4f[T] {
	s := [6]T{
		m[0][0]*m[1][1] - m[1][0]*m[0][1],
		m[0][0]*m[1][2] - m[1][0]*m[0][2],
		m[0][0]*m[1][3] - m[1][0]*m[0][3],
		m[0][1]*m[1][2] - m[1][1]*m[0][2],
		m[0][1]*m[1][3] - m[1][1]*m[0][3],
		m[0][2]*m[1][3] - m[1][2]*m[0][3],
	}
	c := [6]T{
		m[2][0]*m[3][1] - m[3][0]*m[2][1],
		m[2][0]*m[3][2] - m[3][0]*m[2][2],
		m[2][0]*m[3][3] - m[3][0]*m[2][3],
		m[2][1]*m[3][2] - m[3][1]*m[2][2],
		m[2][1]*m[3][3] - m[3][1]*m[2][3],
		m[2][2]*m[3][3] - m[3][2]*m[2][3],
	}

	idet := 1.0 / (s[0]*c[5] - s[1]*c[4] + s[2]*c[3] + s[3]*c[2] - s[4]*c[1] + s[5]*c[0])

	return Matrix4x4f[T]{
		[4]T{
			(m[1][1]*c[5] - m[1][2]*c[4] + m[1][3]*c[3]) * idet,
			(-m[0][1]*c[5] + m[0][2]*c[4] - m[0][3]*c[3]) * idet,
			(m[3][1]*s[5] - m[3][2]*s[4] + m[3][3]*s[3]) * idet,
			(-m[2][1]*s[5] + m[2][2]*s[4] - m[2][3]*s[3]) * idet,
		},
		[4]T{
			(-m[1][0]*c[5] + m[1][2]*c[2] - m[1][3]*c[1]) * idet,
			(m[0][0]*c[5] - m[0][2]*c[2] + m[0][3]*c[1]) * idet,
			(-m[3][0]*s[5] + m[3][2]*s[2] - m[3][3]*s[1]) * idet,
			(m[2][0]*s[5] - m[2][2]*s[2] + m[2][3]*s[1]) * idet,
		},
		[4]T{
			(m[1][0]*c[4] - m[1][1]*c[2] + m[1][3]*c[0]) * idet,
			(-m[0][0]*c[4] + m[0][1]*c[2] - m[0][3]*c[0]) * idet,
			(m[3][0]*s[4] - m[3][1]*s[2] + m[3][3]*s[0]) * idet,
			(-m[2][0]*s[4] + m[2][1]*s[2] - m[2][3]*s[0]) * idet,
		},
		[4]T{
			(-m[1][0]*c[3] + m[1][1]*c[1] - m[1][2]*c[0]) * idet,
			(m[0][0]*c[3] - m[0][1]*c[1] + m[0][2]*c[0]) * idet,
			(-m[3][0]*s[3] + m[3][1]*s[1] - m[3][2]*s[0]) * idet,
			(m[2][0]*s[3] - m[2][1]*s[1] + m[2][2]*s[0]) * idet,
		},
	}
}

func (m Matrix4x4f[T]) ToArrayf32() [4][4]float32 {
	mf := [4][4]float32{}

	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			mf[x][y] = float32(m[x][y])
		}
	}

	return mf
}
