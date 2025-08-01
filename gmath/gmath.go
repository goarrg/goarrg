/*
Copyright 2023 The goARRG Authors.

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

func DegToRad[T constraints.Float](d T) T {
	return d * (math.Pi / 180)
}

func RadToDeg[T constraints.Float](r T) T {
	return r * (180 / math.Pi)
}

func Clamp[N constraints.Float | constraints.Integer](t, low, high N) N {
	return min(max(t, low), high)
}

func InRange[N constraints.Float | constraints.Integer](t, low, high N) bool {
	return (t >= low) && (t <= high)
}
