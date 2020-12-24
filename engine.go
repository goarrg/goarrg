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

package goarrg

import (
	"sync/atomic"
	"time"

	"goarrg.com/debug"
	"goarrg.com/input"
	"goarrg.com/internal/trace"
)

type Platform interface {
	Init() error
	AudioInit(Audio) error
	DisplayInit(Renderer) error
	Update() input.Snapshot
	Shutdown()
	Destroy()
}

type Program interface {
	Init() error
	Update(float64, input.Snapshot)
	Shutdown() bool
	Destroy()
}

type Config struct {
	Platform Platform
	Audio    Audio
	Renderer Renderer
	Program  Program
}

const (
	stateRunning = iota + 1
	stateShutdown
	stateTerminated
)

var system struct {
	platform Platform
	audio    Audio
	renderer Renderer
	program  Program

	state int32
}

func Run(cfg Config) error {
	start := time.Now()
	debug.LogV("Initializing engine")

	defer debug.LogV("Engine terminated")
	defer trace.Shutdown()
	defer atomic.StoreInt32(&system.state, stateTerminated)

	if cfg.Platform == nil {
		return debug.ErrorNew("Invalid platform")
	}

	if cfg.Renderer == nil {
		return debug.ErrorNew("Invalid renderer")
	}

	if cfg.Program == nil {
		return debug.ErrorNew("Invalid program")
	}

	system.platform = cfg.Platform
	system.audio = cfg.Audio
	system.renderer = cfg.Renderer
	system.program = cfg.Program

	if err := system.platform.Init(); err != nil {
		return debug.ErrorWrap(err, "Failed to init platform")
	}

	defer system.platform.Destroy()

	if err := system.platform.DisplayInit(cfg.Renderer); err != nil {
		return debug.ErrorWrap(err, "Failed to init platform display")
	}

	defer system.renderer.Destroy()

	if system.audio != nil {
		if err := system.platform.AudioInit(cfg.Audio); err != nil {
			return debug.ErrorWrap(err, "Failed to init platform audio")
		}

		defer system.audio.Destroy()
	} else {
		system.audio = audioNull{}
	}

	if err := system.program.Init(); err != nil {
		return debug.ErrorWrap(err, "Failed to init user program")
	}

	defer system.program.Destroy()

	atomic.StoreInt32(&system.state, stateRunning)
	debug.LogV("Engine Init took: %v", time.Since(start))

	deltaTime := float64(0.0)

	for Running() {
		traceEnd := debug.TraceStart("Platform Update")
		input := system.platform.Update()
		traceEnd()

		traceEnd = debug.TraceStart("Program Update")
		system.program.Update(deltaTime, input)
		traceEnd()

		debug.Trace("Audio Update", system.audio.Update)

		traceEnd = debug.TraceStart("Renderer Draw")
		deltaTime = system.renderer.Draw()
		traceEnd()
	}

	return nil
}

func Running() bool {
	return atomic.LoadInt32(&system.state) == stateRunning
}

func Shutdown() {
	debug.LogV("Engine signaled to shutdown")
	t := time.AfterFunc(time.Second, func() {
		if atomic.LoadInt32(&system.state) != stateTerminated {
			panic("deadlock")
		}
	})

	if system.program.Shutdown() {
		system.audio.Shutdown()
		system.renderer.Shutdown()
		system.platform.Shutdown()
		atomic.StoreInt32(&system.state, stateShutdown)
		return
	}

	t.Stop()
	debug.LogV("Engine canceled shutdown")
}
