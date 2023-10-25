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

	"golang.org/x/exp/constraints"
)

type Grid[T constraints.Float] struct {
	Pos        Point3f[T]
	Bounds     Bounds3f[T]
	Scale      T
	CellBounds Bounds3i
}

type (
	Gridf32 = Grid[float32]
	Gridf64 = Grid[float64]
)

type GridCell struct {
	X int
	Y int
	Z int
}

func (g *Grid[T]) Init(cellCount Vector3i, cellSize T) {
	g.Scale = cellSize

	if cellCount.X > 0 {
		g.CellBounds.Max.X = cellCount.X - 1
		g.Bounds.Max.X = (T(cellCount.X) / 2) * g.Scale
		g.Bounds.Min.X = -g.Bounds.Max.X
	}

	if cellCount.Y > 0 {
		g.CellBounds.Max.Y = cellCount.Y - 1
		g.Bounds.Max.Y = (T(cellCount.Y) / 2) * g.Scale
		g.Bounds.Min.Y = -g.Bounds.Max.Y
	}

	if cellCount.Z > 0 {
		g.CellBounds.Max.Z = cellCount.Z - 1
		g.Bounds.Max.Z = (T(cellCount.Z) / 2) * g.Scale
		g.Bounds.Min.Z = -g.Bounds.Max.Z
	}
}

func (g *Grid[T]) WorldPosToCell(p Point3f[T]) GridCell {
	x := int(math.Floor(float64((p.X - g.Pos.X - g.Bounds.Min.X) / g.Scale)))
	y := int(math.Floor(float64((p.Y - g.Pos.Y - g.Bounds.Min.Y) / g.Scale)))
	z := int(math.Floor(float64((p.Z - g.Pos.Z - g.Bounds.Min.Z) / g.Scale)))

	return GridCell{
		X: x,
		Y: y,
		Z: z,
	}
}

func (g *Grid[T]) CellToWorldPos(gc GridCell) Point3f[T] {
	return Point3f[T]{
		X: (T(gc.X) * g.Scale) + g.Bounds.Min.X + (0.5 * g.Scale) + g.Pos.X,
		Y: (T(gc.Y) * g.Scale) + g.Bounds.Min.Y + (0.5 * g.Scale) + g.Pos.Y,
		Z: (T(gc.Z) * g.Scale) + g.Bounds.Min.Z + (0.5 * g.Scale) + g.Pos.Z,
	}
}

func (gc GridCell) Clamp(bounds Bounds3i) GridCell {
	if gc.X < bounds.Min.X {
		gc.X = bounds.Min.X
	}

	if gc.X > bounds.Max.X {
		gc.X = bounds.Max.X
	}

	if gc.Y < bounds.Min.Y {
		gc.Y = bounds.Min.Y
	}

	if gc.Y > bounds.Max.Y {
		gc.Y = bounds.Max.Y
	}

	if gc.Z < bounds.Min.Z {
		gc.Z = bounds.Min.Z
	}

	if gc.Z > bounds.Max.Z {
		gc.Z = bounds.Max.Z
	}

	return gc
}

func (gc GridCell) ToArray() [3]int {
	return [3]int{
		gc.X,
		gc.Y,
		gc.Z,
	}
}
