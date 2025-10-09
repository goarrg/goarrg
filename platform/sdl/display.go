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

package sdl

import (
	"goarrg.com"
	"goarrg.com/debug"
	"goarrg.com/gmath"
)

type displaySystem struct {
	mainWindow *window
}

func (*platform) DisplayInit(renderer goarrg.Renderer) error {
	if vk, ok := renderer.(goarrg.VkRenderer); ok {
		if err := vkInit(vk); err != nil {
			if gl, ok := renderer.(goarrg.GLRenderer); ok {
				Platform.logger.IPrintf("Failed to init vk renderer %v", err)
				Platform.logger.IPrintf("Falling back to gl renderer")

				return glInit(gl)
			}
			return err
		}

		return nil
	}

	if gl, ok := renderer.(goarrg.GLRenderer); ok {
		if err := glInit(gl); err != nil {
			return err
		}

		return nil
	}

	return debug.Errorf("Invalid renderer")
}

func (d *displaySystem) hasKeyboardFocus() bool {
	return d.mainWindow.keyboardFocus
}

func (d *displaySystem) hasMouseFocus() bool {
	return d.mainWindow.mouseFocus
}

func (d *displaySystem) pointInsideWindow(p gmath.Point3f64) bool {
	return d.mainWindow.windowExtent.CheckPoint(p.Subtract(gmath.Vector3f64(d.mainWindow.windowPos)))
}

func (d *displaySystem) globalPointToRelativePoint(p gmath.Point3f64) gmath.Point3f64 {
	p = d.mainWindow.windowExtent.ClampPoint(p.Subtract(gmath.Vector3f64(d.mainWindow.windowPos)))
	if d.mainWindow.windowExtent == d.mainWindow.surfaceExtent {
		return p
	}

	return gmath.Point3f64(
		gmath.Vector3f64(p).
			ScaleInverse(gmath.Vector3f64(d.mainWindow.windowExtent)).
			Scale(gmath.Vector3f64(d.mainWindow.surfaceExtent)),
	)
}

func (d *displaySystem) destroy() {
	if d.mainWindow != nil {
		d.mainWindow.destroy()
	}
}
