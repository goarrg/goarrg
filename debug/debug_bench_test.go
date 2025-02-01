//go:build !goarrg_build_debug
// +build !goarrg_build_debug

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
	"testing"
)

type nilWritter struct{}

func (*nilWritter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func BenchmarkLog(b *testing.B) {
	loggerOut = &nilWritter{}
	logger := NewLogger()
	logger.SetLevel(LogLevelVerbose)

	for n := 0; n < b.N; n++ {
		logger.VPrint("TEST")
	}
}

func BenchmarkLogIgnore(b *testing.B) {
	loggerOut = &nilWritter{}
	logger := NewLogger()
	logger.SetLevel(LogLevelError)

	for n := 0; n < b.N; n++ {
		logger.VPrint("TEST")
	}
}

func BenchmarkError(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = Errorf("TEST")
		}
	})

	err := Errorf("TEST")
	_ = fmt.Sprintf("%v", err)

	b.Run("Printf", func(b *testing.B) {
		b.Run("s", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = fmt.Sprintf("%s", err)
			}
		})
		b.Run("v", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = fmt.Sprintf("%v", err)
			}
		})
	})

	b.Run("Wrap", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = ErrorWrapf(err, "TEST")
		}
	})

	err = ErrorWrapf(err, "TEST")
	_ = fmt.Sprintf("%v", err)

	b.Run("PrintfWrapped", func(b *testing.B) {
		b.Run("s", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = fmt.Sprintf("%s", err)
			}
		})
		b.Run("v", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = fmt.Sprintf("%v", err)
			}
		})
	})
}

func BenchmarkStackTrace(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = StackTrace(0)
	}
}

func stackStringTest() string {
	return StackTrace(0)
}

func nestedStackStringTest() string {
	return stackStringTest()
}

func nestedNestedStackStringTest() string {
	return nestedStackStringTest()
}

func BenchmarkNestedStackTrace(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = nestedNestedStackStringTest()
	}
}
