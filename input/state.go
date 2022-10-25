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

package input

import "goarrg.com/gmath"

type State interface {
	isInputState()
}

type StateDelta interface {
	isInputStateDelta()
}

/*
Value represents a linear action such as a button or analog trigger.
For buttons this is either a 0 or a 1 or a delta of -1.
For analog actions this would be between 0.0 and 1.0 inclusive.
*/
type Value float32

/*
Axis represents a 3D vector action such as mouse or analog stick movement.
*/
type Axis struct {
	gmath.Vector3f64
}

/*
Coords represents a 3D position action such as mouse or tap position.
*/
type Coords struct {
	gmath.Point3f64
}

func (v Value) isInputState() {
}

func (v Value) isInputStateDelta() {
}

func (a Axis) isInputState() {
}

func (a Axis) isInputStateDelta() {
}

func (a Coords) isInputState() {
}
