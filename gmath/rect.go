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

type Rect[T constraints.Integer | constraints.Float] struct {
	X, Y, W, H T
}

type (
	Rectint = Rect[int]
	Recti32 = Rect[int32]
	Recti64 = Rect[int64]

	Rectf32 = Rect[float32]
	Rectf64 = Rect[float64]
)
