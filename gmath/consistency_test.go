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
	"reflect"
	"strings"
	"testing"
)

func TestConsistency(t *testing.T) {
	testForMethods := func(m map[string]reflect.Type, l []string, veto func(string) bool) {
		for _, f := range l {
			for k, v := range m {
				if veto(k) {
					continue
				}
				_, ok := v.MethodByName(f)
				if !ok {
					t.Fatal(k, "does not implement", f)
				}
			}
		}
	}
	{
		vectorTypeMap := map[string]reflect.Type{
			"Point2i": reflect.TypeFor[Point2int](),
			"Point2f": reflect.TypeFor[Point2f32](),

			"Point3i": reflect.TypeFor[Point3int](),
			"Point3f": reflect.TypeFor[Point3f32](),

			"Vector2i": reflect.TypeFor[Vector2int](),
			"Vector2f": reflect.TypeFor[Vector2f32](),

			"Vector3i": reflect.TypeFor[Vector3int](),
			"Vector3f": reflect.TypeFor[Vector3f32](),
		}
		testForMethods(vectorTypeMap,
			[]string{
				"Add",
				"Subtract",
				"ToArray",
			},
			func(s string) bool { return false },
		)
		testForMethods(vectorTypeMap,
			[]string{
				"IsNAN",
				"ToArrayf32",
			},
			func(s string) bool { return !strings.HasSuffix(s, "f") },
		)
		testForMethods(vectorTypeMap,
			[]string{
				"ToArrayi32",
			},
			func(s string) bool { return !strings.HasSuffix(s, "i") },
		)
		testForMethods(vectorTypeMap,
			[]string{
				"VectorTo",
			},
			func(s string) bool { return !strings.HasPrefix(s, "Point") },
		)
		testForMethods(vectorTypeMap,
			[]string{
				"Abs",
				"Min",
				"Max",
				"Clamp",
				"Add",
				"AddUniform",
				"Subtract",
				"SubtractUniform",
				"Scale",
				"ScaleInverse",
				"ScaleUniform",
				"ScaleInverseUniform",
			},
			func(s string) bool { return !strings.HasPrefix(s, "Vector") },
		)
		testForMethods(vectorTypeMap,
			[]string{
				"Dot",
				"Magnitude",
				"Angle",
				"Normalize",
			},
			func(s string) bool { return !strings.HasPrefix(s, "Vector") || !strings.HasSuffix(s, "f") },
		)
	}
	{
		rangeTypeMap := map[string]reflect.Type{
			"Extent2i": reflect.TypeFor[Extent2int](),
			"Extent2f": reflect.TypeFor[Extent2f32](),

			"Extent3i": reflect.TypeFor[Extent3int](),
			"Extent3f": reflect.TypeFor[Extent3f32](),

			"Bounds2i": reflect.TypeFor[Bounds2int](),
			"Bounds2f": reflect.TypeFor[Bounds2f32](),

			"Bounds3i": reflect.TypeFor[Bounds3int](),
			"Bounds3f": reflect.TypeFor[Bounds3f32](),
		}
		testForMethods(rangeTypeMap,
			[]string{
				"CheckPoint",
				"CheckVector",
				"ClampPoint",
				"ClampVector",
			},
			func(s string) bool { return false },
		)
		testForMethods(rangeTypeMap,
			[]string{
				"InRange",
				"Min",
				"Max",
				"Clamp",
			},
			func(s string) bool { return !strings.HasPrefix(s, "Extent") },
		)
		testForMethods(rangeTypeMap,
			[]string{
				"Area",
			},
			func(s string) bool { return !strings.HasPrefix(s, "Extent2") },
		)
		testForMethods(rangeTypeMap,
			[]string{
				"Volume",
			},
			func(s string) bool { return !strings.HasPrefix(s, "Extent3") },
		)
	}
}
