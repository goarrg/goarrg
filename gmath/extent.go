/*
Copyright 2024 The goARRG Authors.

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

type Extent3[T constraints.Integer | constraints.Float] struct {
	X, Y, Z T
}

type (
	Extent3f32 = Extent3[float32]
	Extent3f64 = Extent3[float64]

	Extent3int = Extent3[int]
	Extent3i32 = Extent3[int32]
	Extent3i64 = Extent3[int64]

	Extent3uint = Extent3[uint]
	Extent3u32  = Extent3[uint32]
	Extent3u64  = Extent3[uint64]
)

func (e Extent3[T]) InRange(min, max Extent3[T]) bool {
	return InRange(e.X, min.X, max.X) &&
		InRange(e.Y, min.Y, max.Y) &&
		InRange(e.Z, min.Z, max.Z)
}

func (e Extent3[T]) Volume() T {
	return e.X * e.Y * e.Z
}
