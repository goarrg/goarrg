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
)

type Grid struct {
	Pos        Point3f64
	Bounds     Bounds3f64
	Scale      float64
	CellBounds Bounds3i
}

type GridCell struct {
	X int
	Y int
	Z int
}

func (g *Grid) Init(cellCount Vector3i, cellSize float64) {
	g.Scale = cellSize

	if cellCount.X > 0 {
		g.CellBounds.Max.X = cellCount.X - 1
		g.Bounds.Max.X = (float64(cellCount.X) / 2) * g.Scale
		g.Bounds.Min.X = -g.Bounds.Max.X
	}

	if cellCount.Y > 0 {
		g.CellBounds.Max.Y = cellCount.Y - 1
		g.Bounds.Max.Y = (float64(cellCount.Y) / 2) * g.Scale
		g.Bounds.Min.Y = -g.Bounds.Max.Y
	}

	if cellCount.Z > 0 {
		g.CellBounds.Max.Z = cellCount.Z - 1
		g.Bounds.Max.Z = (float64(cellCount.Z) / 2) * g.Scale
		g.Bounds.Min.Z = -g.Bounds.Max.Z
	}
}

func (g *Grid) WorldPosToCell(p Point3f64) GridCell {
	x := int(math.Floor((p.X - g.Pos.X - g.Bounds.Min.X) / g.Scale))
	y := int(math.Floor((p.Y - g.Pos.Y - g.Bounds.Min.Y) / g.Scale))
	z := int(math.Floor((p.Z - g.Pos.Z - g.Bounds.Min.Z) / g.Scale))

	return GridCell{
		X: x,
		Y: y,
		Z: z,
	}
}

func (g *Grid) CellToWorldPos(gc GridCell) Point3f64 {
	return Point3f64{
		X: (float64(gc.X) * g.Scale) + g.Bounds.Min.X + (0.5 * g.Scale) + g.Pos.X,
		Y: (float64(gc.Y) * g.Scale) + g.Bounds.Min.Y + (0.5 * g.Scale) + g.Pos.Y,
		Z: (float64(gc.Z) * g.Scale) + g.Bounds.Min.Z + (0.5 * g.Scale) + g.Pos.Z,
	}
}

func (gc GridCell) Clamp(g Grid) GridCell {
	if gc.X < g.CellBounds.Min.X {
		gc.X = g.CellBounds.Min.X
	}

	if gc.X > g.CellBounds.Max.X {
		gc.X = g.CellBounds.Max.X
	}

	if gc.Y < g.CellBounds.Min.Y {
		gc.Y = g.CellBounds.Min.Y
	}

	if gc.Y > g.CellBounds.Max.Y {
		gc.Y = g.CellBounds.Max.Y
	}

	if gc.Z < g.CellBounds.Min.Z {
		gc.Z = g.CellBounds.Min.Z
	}

	if gc.Z > g.CellBounds.Max.Z {
		gc.Z = g.CellBounds.Max.Z
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
