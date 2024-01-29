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
	"goarrg.com/debug"
	"goarrg.com/gmath"
)

type WindowMode int

const (
	WindowModeWindowed WindowMode = iota
	WindowModeBorderless
	WindowModeFullscreen
)

type AudioImporterConfig struct {
	EnableWAV bool
}

type AudioConfig struct {
	Importer AudioImporterConfig
}

type WindowConfig struct {
	Title string
	Rect  gmath.Rectint
	Mode  WindowMode
}

type Config struct {
	Audio  AudioConfig
	Window WindowConfig
}

func Setup(cfg Config) error {
	{
		if cfg.Window.Title == "" {
			cfg.Window.Title = "goarrg SDL"
		}

		if cfg.Window.Rect.W <= 0 {
			return debug.Errorf("Invalid window size: %+v", cfg.Window.Rect)
		}

		if cfg.Window.Rect.H <= 0 {
			return debug.Errorf("Invalid window size: %+v", cfg.Window.Rect)
		}

		if cfg.Window.Mode > WindowModeFullscreen {
			return debug.Errorf("Invalid window mode: %d", cfg.Window.Mode)
		}
	}

	Platform.config = cfg

	return nil
}
