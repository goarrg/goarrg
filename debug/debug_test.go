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
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func stackTest() error {
	return Errorf("Stack test")
}

func passthroughTest() error {
	return stackTest()
}

func nestedStackTest() error {
	return ErrorWrapf(ErrorWrapf(passthroughTest(), "Nested stack test"), "Double nested test")
}

func nestedNestedStackTest() error {
	return ErrorWrapf(nestedStackTest(), "Nested nested stack test")
}

func checkOutput(t *testing.T, b *strings.Builder, enabled bool, want string) {
	t.Helper()

	if enabled {
		r := regexp.MustCompile(strings.TrimSpace(want))

		if r.MatchString(b.String()) {
			t.Log(strings.TrimSpace(b.String()))
		} else {
			t.Logf("\nExpected: %s\nGot: %s", strings.TrimSpace(want), b.String())
			t.FailNow()
		}

		b.Reset()
	} else {
		if b.String() != "" {
			t.Logf("\nExpected: \nGot: %s", b.String())
			t.FailNow()
		}
	}
}

func runLoggerTests(t *testing.T, b *strings.Builder, level uint32, l *Logger, tags ...string) {
	tagString := ""

	for _, t := range tags {
		tagString += `\[` + t + `\] `
	}

	l.VPrint(0, "VPrint Test", 1)
	checkOutput(t, b, LogLevelVerbose >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[VERBOSE\] %s0VPrint Test1`, tagString))

	l.VPrintln(0, "VPrintln Test", 1)
	checkOutput(t, b, LogLevelVerbose >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[VERBOSE\] %s0 VPrintln Test 1`, tagString))

	l.VPrintf("VPrintf Test %d", 123)
	checkOutput(t, b, LogLevelVerbose >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[VERBOSE\] %sVPrintf Test 123`, tagString))

	l.IPrint(0, "IPrint Test", 1)
	checkOutput(t, b, LogLevelInfo >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[INFO\]    %s0IPrint Test1`, tagString))

	l.IPrintln(0, "IPrintln Test", 1)
	checkOutput(t, b, LogLevelInfo >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[INFO\]    %s0 IPrintln Test 1`, tagString))

	l.IPrintf("IPrintf Test %d", 123)
	checkOutput(t, b, LogLevelInfo >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[INFO\]    %sIPrintf Test 123`, tagString))

	l.WPrint(0, "WPrint Test", 1)
	checkOutput(t, b, LogLevelWarn >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[WARN\]    %s0WPrint Test1`, tagString))

	l.WPrintln(0, "WPrintln Test", 1)
	checkOutput(t, b, LogLevelWarn >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[WARN\]    %s0 WPrintln Test 1`, tagString))

	l.WPrintf("WPrintf Test %d", 123)
	checkOutput(t, b, LogLevelWarn >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[WARN\]    %sWPrintf Test 123`, tagString))

	l.EPrint(0, "EPrint Test", 1)
	checkOutput(t, b, LogLevelError >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[ERROR\]   %s0EPrint Test1`, tagString))

	l.EPrintln(0, "EPrintln Test", 1)
	checkOutput(t, b, LogLevelError >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[ERROR\]   %s0 EPrintln Test 1`, tagString))

	l.EPrintf("EPrintf Test %d", 123)
	checkOutput(t, b, LogLevelError >= level, fmt.Sprintf(`\d*:\d*:\d*.\d* \[ERROR\]   %sEPrintf Test 123`, tagString))
}

func Test_debug(t *testing.T) {
	b := &strings.Builder{}
	loggerOut = b

	// tests global log level
	for level := LogLevelVerbose; level <= LogLevelError; level++ {
		t.Logf("Testing global log level: %d", level)
		SetLevel(level)

		// global logger tests
		t.Logf("Testing global logger")
		runLoggerTests(t, b, level, &logger)

		standardLogger := NewLogger()
		t.Logf("Testing standard logger")
		runLoggerTests(t, b, level, standardLogger)

		taggedLogger := NewLogger("TAG")
		t.Logf("Testing tag logger")
		runLoggerTests(t, b, level, taggedLogger, "TAG")

		taggedLogger = taggedLogger.NewLoggerWithTags("TAG2")
		t.Logf("Testing nested tag logger")
		runLoggerTests(t, b, level, taggedLogger, "TAG", "TAG2")
	}

	// test individual log level
	for level := LogLevelVerbose; level <= LogLevelError; level++ {
		t.Logf("Testing individual log level: %d", level)

		// global logger tests which should be at LogLevelError from our
		// previous tests
		t.Logf("Testing global logger")
		runLoggerTests(t, b, LogLevelError, &logger)

		standardLogger := NewLogger()
		standardLogger.SetLevel(level)
		t.Logf("Testing standard logger")
		runLoggerTests(t, b, level, standardLogger)

		taggedLogger := NewLogger("TAG")
		taggedLogger.SetLevel(level)
		t.Logf("Testing tag logger")
		runLoggerTests(t, b, level, taggedLogger, "TAG")

		taggedLogger = taggedLogger.NewLoggerWithTags("TAG2")
		taggedLogger.SetLevel(level)
		t.Logf("Testing nested tag logger")
		runLoggerTests(t, b, level, taggedLogger, "TAG", "TAG2")
	}

	err := Errorf("Err Test")

	EPrintf("%s", err)
	checkOutput(t, b, true, `\d*:\d*:\d*.\d* \[ERROR\]   Err Test`)

	EPrint(err)
	checkOutput(t, b, true, `
\d*:\d*:\d*.\d* \[ERROR\]   Err Test
	.*
		.*\S:\d*
`)

	err = ErrorWrapf(err, "Chain test")

	EPrintf("%s", err)
	checkOutput(t, b, true, `\d*:\d*:\d*.\d* \[ERROR\]   Chain test: Err Test`)

	EPrint(err)
	checkOutput(t, b, true, `
\d*:\d*:\d*.\d* \[ERROR\]   Err Test
.*
	.*\S:\d*

Chain test
.*
	.*\S:\d*
`)

	err = ErrorWrapf(errors.New("Unknown error"), "Unknown error chain test")

	EPrintf("%s", err)
	checkOutput(t, b, true, `\d*:\d*:\d*.\d* \[ERROR\]   Unknown error chain test: Unknown error`)

	EPrint(err)
	checkOutput(t, b, true, `
\d*:\d*:\d*.\d* \[ERROR\]   Unknown error

Unknown error chain test
	.*
		.*\S:\d*
`)

	EPrint(stackTest())
	checkOutput(t, b, true, `
\d*:\d*:\d*.\d* \[ERROR\]   Stack test
.*
	.*\S:\d*
.*
	.*\S:\d*
`)

	EPrint(nestedNestedStackTest())
	checkOutput(t, b, true, `
\d*:\d*:\d*.\d* \[ERROR\]   Stack test
.*
	.*\S:\d*
.*
	.*\S:\d*

Nested stack test
.*
	.*\S:\d*

Double nested test
.*
	.*\S:\d*

Nested nested stack test
.*
	.*\S:\d*
.*
	.*\S:\d*
`)

	EPrintf("%s", StackTrace(0))
	checkOutput(t, b, true, `
\d*:\d*:\d*.\d* \[ERROR\]   .*
	.*\S:\d*
`)
}
