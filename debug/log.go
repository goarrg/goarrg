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
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

const (
	LogLevelVerbose = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var (
	logOut   = log.New(os.Stderr, "", log.Lmicroseconds).Output
	logLevel uint32
)

func LogSetLevel(l uint32) {
	if l <= LogLevelError {
		atomic.StoreUint32(&logLevel, l)
	}
}

func LogWillLog(l uint32) bool {
	return atomic.LoadUint32(&logLevel) <= l
}

// LogV is for high detail information spam logging for when nothing else is giving enough information
func LogV(format string, args ...interface{}) {
	if atomic.LoadUint32(&logLevel) > LogLevelVerbose {
		return
	}

	_ = logOut(0, "[VERBOSE] "+strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// LogI is for basic and succinct high level logs that may be useful for general debugging
func LogI(format string, args ...interface{}) {
	if atomic.LoadUint32(&logLevel) > LogLevelInfo {
		return
	}

	_ = logOut(0, "[INFO]    "+strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// LogW is for things that generally won't cause a problem, but a good place to look if there is a problem
func LogW(format string, args ...interface{}) {
	if atomic.LoadUint32(&logLevel) > LogLevelWarn {
		return
	}

	_ = logOut(0, "[WARN]    "+strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// LogE is for clear problems that the dev should/must fix
func LogE(format string, args ...interface{}) {
	_ = logOut(0, "[ERROR]   "+strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// LogErr is a helper method for error logging and checking
func LogErr(err error) bool {
	if err == nil {
		return false
	}

	_ = logOut(0, "[ERROR]   "+fmt.Sprintf("%v", err))
	return true
}
