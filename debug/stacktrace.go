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

package debug

import (
	"runtime"
	"strconv"
	"strings"
)

type stack struct {
	str string
	pcs []uintptr
}

func (s *stack) String() string {
	if s.str == "" {
		s.str = stackString(s.pcs, 1)
	}

	return s.str
}

func callers(skip int) []uintptr {
	buf := make([]uintptr, 8)

	for {
		n := runtime.Callers(2+skip, buf)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]uintptr, len(buf)+8)
	}
}

func stackString(callers []uintptr, nTabs int) string {
	tabs := strings.Repeat("\t", nTabs)
	frames := runtime.CallersFrames(callers)
	stack := ""
	frame, more := frames.Next()

	switch {
	case strings.HasPrefix(frame.Function, "runtime."):
	case strings.HasPrefix(frame.Function, "testing."):
	case strings.HasPrefix(frame.Function[strings.IndexAny(frame.Function, ".")+1:], "_cgo"):
	default:
		stack += frame.Function + "\n" + tabs + "\t" + frame.File + ":" + strconv.Itoa(frame.Line)
	}

	for more {
		frame, more = frames.Next()

		switch {
		case strings.HasPrefix(frame.Function, "runtime."):
			continue
		case strings.HasPrefix(frame.Function, "testing."):
			continue
		case strings.HasPrefix(frame.Function[strings.IndexAny(frame.Function, ".")+1:], "_cgo"):
			continue
		}

		stack += "\n" + tabs + frame.Function + "\n" + tabs + "\t" + frame.File + ":" + strconv.Itoa(frame.Line)
	}

	return strings.TrimSpace(stack)
}

func StackTrace(skip int) string {
	return stackString(callers(skip+1), 0)
}
