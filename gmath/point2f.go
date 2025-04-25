/*
Copyright 2025 The goARRG Authors.

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

	"golang.org/x/exp/constraints"
)

type Point2f[T constraints.Float] struct {
	X, Y T
}

type (
	Point2f32 = Point2f[float32]
	Point2f64 = Point2f[float64]
)

func Point2fFromArray[T constraints.Float](a [2]T) Point2f[T] {
	return Point2f[T]{
		X: a[0],
		Y: a[1],
	}
}

func (p Point2f[T]) VectorTo(p2 Point2f[T]) Vector2f[T] {
	return Vector2f[T]{
		X: p2.X - p.X,
		Y: p2.Y - p.Y,
	}
}

func (p Point2f[T]) Add(v Vector2f[T]) Point2f[T] {
	return Point2f[T]{
		X: p.X + v.X,
		Y: p.Y + v.Y,
	}
}

func (p Point2f[T]) Subtract(v Vector2f[T]) Point2f[T] {
	return Point2f[T]{
		X: p.X - v.X,
		Y: p.Y - v.Y,
	}
}

func (p Point2f[T]) IsNAN() bool {
	return math.IsNaN(float64(p.X)) || math.IsNaN(float64(p.Y))
}

func (p Point2f[T]) ToArray() [2]T {
	return [2]T{p.X, p.Y}
}

func (p Point2f[T]) ToArrayf32() [2]float32 {
	return [2]float32{float32(p.X), float32(p.Y)}
}
