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

import "golang.org/x/exp/constraints"

type Point2i[T constraints.Integer] struct {
	X, Y T
}

type (
	Point2int = Point2i[int]
	Point2i32 = Point2i[int32]
	Point2i64 = Point2i[int64]
)

func Point2iFromArray[T constraints.Integer](a [2]T) Point2i[T] {
	return Point2i[T]{
		X: a[0],
		Y: a[1],
	}
}

func (p Point2i[T]) VectorTo(p2 Point2i[T]) Vector2i[T] {
	return Vector2i[T]{
		X: p2.X - p.X,
		Y: p2.Y - p.Y,
	}
}

func (p Point2i[T]) Add(v Vector2i[T]) Point2i[T] {
	return Point2i[T]{
		X: p.X + v.X,
		Y: p.Y + v.Y,
	}
}

func (p Point2i[T]) Subtract(v Vector2i[T]) Point2i[T] {
	return Point2i[T]{
		X: p.X - v.X,
		Y: p.Y - v.Y,
	}
}

func (p Point2i[T]) ToArray() [2]T {
	return [2]T{p.X, p.Y}
}

func (p Point2i[T]) ToArrayi32() [2]int32 {
	return [2]int32{int32(p.X), int32(p.Y)}
}
