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
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"goarrg.com/debug"
	"goarrg.com/internal/trace"
)

type Platform interface {
	Init() error
	AudioInit(Audio) error
	DisplayInit(Renderer) error
	Update()
	Destroy()
}

type Program interface {
	Init() error
	Update(float64)
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
	stateShutdownConfirmed
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

	// setup signal handlers to force shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer close(c)
	defer signal.Stop(c)

	go func() {
		if <-c == nil {
			debug.LogV("Engine closed signal handler")
			return
		}

		debug.LogV("Engine signaled to force shutdown")
		atomic.StoreInt32(&system.state, stateShutdownConfirmed)

		time.AfterFunc(time.Second, func() {
			if atomic.LoadInt32(&system.state) != stateTerminated {
				panic("deadlock")
			}
		})
	}()

	atomic.StoreInt32(&system.state, stateRunning)
	debug.LogV("Engine Init took: %v", time.Since(start))

	deltaTime := float64(0.0)

loop:

	for atomic.LoadInt32(&system.state) == stateRunning {
		debug.Trace("Platform Update", system.platform.Update)

		traceEnd := debug.TraceStart("Program Update")
		system.program.Update(deltaTime)
		traceEnd()

		debug.Trace("Audio Update", system.audio.Update)

		traceEnd = debug.TraceStart("Renderer Draw")
		deltaTime = system.renderer.Draw()
		traceEnd()
	}

	if atomic.LoadInt32(&system.state) == stateShutdownConfirmed {
		return nil
	}

	t := time.AfterFunc(time.Second, func() {
		if atomic.LoadInt32(&system.state) != stateTerminated {
			panic("deadlock")
		}
	})

	if !system.program.Shutdown() {
		t.Stop()
		atomic.StoreInt32(&system.state, stateRunning)
		debug.LogV("Engine canceled shutdown")
		goto loop
	}

	atomic.StoreInt32(&system.state, stateShutdownConfirmed)

	return nil
}

/*
	Running reports whether the engine is considered to still be running.
	It will only return true when the main loop is running and before shutdown
	confirmation from Program.Shutdown(). This is so that you can depend on
	Running() to terminate your loops/threads.

	SIGINT will bypass Program.Shutdown() and force a terminate. This is so we
	have a easy way to terminate the engine in the event of deadlocks.
*/
func Running() bool {
	s := atomic.LoadInt32(&system.state)
	return s == stateRunning || s == stateShutdown
}

/*
	Shutdown is a thread safe signal to the engine that it should shutdown.
	The signal usually would come from the Platform or Program packages.
*/
func Shutdown() {
	if atomic.CompareAndSwapInt32(&system.state, stateRunning, stateShutdown) {
		debug.LogV("Engine signaled to shutdown")
	}
}
