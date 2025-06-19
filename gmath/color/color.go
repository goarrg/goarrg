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

package color

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

type convertable interface {
	toUNorm64() UNorm[float64]
	newFromUNorm64(UNorm[float64]) any
}

type Convertable[T constraints.Unsigned | constraints.Float] interface {
	convertable
	~struct{ R, G, B, A T }
}

func normalize[T constraints.Unsigned | constraints.Float](in T) float64 {
	switch any(in).(type) {
	case uint8:
		return float64(in) / float64(math.MaxUint8)
	case uint16:
		return float64(in) / float64(math.MaxUint16)
	case uint32:
		return float64(in) / float64(math.MaxUint32)
	case uint64:
		return float64(in) / float64(math.MaxUint64)
	case uint:
		return float64(in) / float64(math.MaxUint)
	case float32, float64:
		return float64(in)

	default:
		panic(fmt.Sprintf("Unknown type: %T", in))
	}
}

func unnormalize[T constraints.Unsigned | constraints.Float](in float64) T {
	switch any(T(0)).(type) {
	case uint8:
		return T(float64(in) * float64(math.MaxUint8))
	case uint16:
		return T(float64(in) * float64(math.MaxUint16))
	case uint32:
		return T(float64(in) * float64(math.MaxUint32))
	case uint64:
		return T(float64(in) * float64(math.MaxUint64))
	case uint:
		return T(float64(in) * float64(math.MaxUint))
	case float32, float64:
		return T(in)

	default:
		panic(fmt.Sprintf("Unknown type: %T", in))
	}
}

/*
Convert will convert between the different color formats and component sizes, do not convert non color data such as normals to/from
non linear formats such as SRGB.
*/
func Convert[Tout Convertable[U], Tin Convertable[V], U constraints.Unsigned | constraints.Float, V constraints.Unsigned | constraints.Float](in Tin) Tout {
	return Tout.newFromUNorm64(Tout{}, in.toUNorm64()).(Tout)
}

/*
UNorm represents a 4 component linear color. If you have gotten the color
from an image, color picker or image editor, the color is almost surely SRGB
not linear.
*/
type UNorm[T constraints.Unsigned | constraints.Float] struct {
	R, G, B, A T
}

var _ convertable = UNorm[float32]{}

func (c UNorm[T]) toUNorm64() UNorm[float64] {
	return UNorm[float64]{
		R: normalize(c.R),
		G: normalize(c.G),
		B: normalize(c.B),
		A: normalize(c.A),
	}
}

func (c UNorm[T]) newFromUNorm64(in UNorm[float64]) any {
	return UNorm[T]{
		R: unnormalize[T](in.R),
		G: unnormalize[T](in.G),
		B: unnormalize[T](in.B),
		A: unnormalize[T](in.A),
	}
}

/*
SRGB represents a 4 component SRGB color. If uploading data to the gpu,
ensure that you convert to UNorm or use a SRGB image format or the colors will
be wrong.
*/
type SRGB[T constraints.Unsigned | constraints.Float] struct {
	R, G, B, A T
}

var _ convertable = SRGB[float32]{}

func (c SRGB[T]) toUNorm64() UNorm[float64] {
	f := func(in T) float64 {
		out := normalize(in)
		if out > 0.04045 {
			out = math.Pow((out+0.055)/1.055, 2.4)
		} else {
			out = out / 12.92
		}
		return out
	}
	return UNorm[float64]{
		R: f(c.R),
		G: f(c.G),
		B: f(c.B),
		A: normalize(c.A),
	}
}

func (c SRGB[T]) newFromUNorm64(in UNorm[float64]) any {
	f := func(in float64) T {
		if in > 0.0031308 {
			in = (1.055 * math.Pow(in, 1/2.4)) - 0.055
		} else {
			in = in * 12.92
		}
		return unnormalize[T](in)
	}
	c.R = f(in.R)
	c.G = f(in.G)
	c.B = f(in.B)
	c.A = unnormalize[T](in.A)
	return c
}
