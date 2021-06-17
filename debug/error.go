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
	"errors"
	"fmt"
	"runtime"
)

type errorNode struct {
	stack stack
	msg   string
	next  error
}

func Errorf(format string, args ...interface{}) error {
	if len(args) == 0 {
		return &errorNode{
			stack: stack{pcs: callers(1)},
			msg:   format,
		}
	}

	return &errorNode{
		stack: stack{pcs: callers(1)},
		msg:   fmt.Sprintf(format, args...),
	}
}

func ErrorWrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	var pcs []uintptr

	if en, ok := err.(*errorNode); ok {
		sameStack := false
		pcs = make([]uintptr, 2)

		if runtime.Callers(2, pcs) != 2 {
			panic("WTF")
		}

		for i, pc := range en.stack.pcs {
			if pc == pcs[1] {
				pcs = append(pcs, en.stack.pcs[i+1:]...)

				if i > 1 {
					en.stack.pcs = en.stack.pcs[:i-1]
				} else {
					en.stack.pcs = en.stack.pcs[:1]
				}

				sameStack = true
				break
			}
		}

		if !sameStack {
			pcs = callers(1)
		}
	} else {
		pcs = callers(1)
	}

	if len(args) == 0 {
		return &errorNode{
			stack: stack{pcs: pcs},
			msg:   format,
			next:  err,
		}
	}

	return &errorNode{
		stack: stack{pcs: pcs},
		msg:   fmt.Sprintf(format, args...),
		next:  err,
	}
}

func (en *errorNode) Error() string {
	if en.next != nil {
		return en.msg + ": " + en.next.Error()
	}

	return en.msg
}

func (en *errorNode) Unwrap() error {
	return en.next
}

func (en *errorNode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if next := errors.Unwrap(en); next != nil {
			fmt.Fprintf(s, "%v", next)
			s.Write([]byte("\n\n"))
		}

		s.Write([]byte(en.msg + "\n\t" + en.stack.String()))

	case 's':
		s.Write([]byte(en.Error()))
	}
}
