//go:build debug
// +build debug

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

package trace

import (
	"context"
	"os"
	"runtime/trace"
)

var (
	traceCtx  context.Context
	traceTask *trace.Task
	out       *os.File
)

func init() {
	if !trace.IsEnabled() {
		file, err := os.Create("out.trace")
		if err != nil {
			panic(err)
		}

		out = file
		err = trace.Start(file)

		if err != nil {
			panic(err)
		}
	}

	traceCtx, traceTask = trace.NewTask(context.Background(), "Debug")
}

func Do(name string, f func()) {
	trace.WithRegion(traceCtx, name, f)
}

func Start(name string) func() {
	return trace.StartRegion(traceCtx, name).End
}

func Shutdown() {
	traceTask.End()
	trace.Stop()
	out.Close()
}
