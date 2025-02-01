//go:build goarrg_build_debug
// +build goarrg_build_debug

/*
Copyright 2023 The goARRG Authors.

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
	"syscall"
)

func initTerminal() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	p := kernel32.NewProc("AttachConsole")
	ret, _, _ := syscall.SyscallN(p.Addr(), ^uintptr(0))

	// if we can't attach then alloc, but make sure there isn't a console by calling free first
	if ret == 0 {
		p = kernel32.NewProc("FreeConsole")
		_, _, _ = syscall.SyscallN(p.Addr())
		p = kernel32.NewProc("AllocConsole")
		ret, _, _ = syscall.SyscallN(p.Addr())
	}

	return ret != 0
}
