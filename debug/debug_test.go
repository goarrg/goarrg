//+build !debug

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

func runLoggerTests(l *Logger) {
	l.SetLevel(LogLevelVerbose)
	l.VPrint(0, "VPrint Test", 1)
	l.VPrintln(0, "VPrintln Test", 1)
	l.VPrintf("VPrintf Test %d", 123)

	l.IPrint(0, "IPrint Test", 1)
	l.IPrintln(0, "IPrintln Test", 1)
	l.IPrintf("IPrintf Test %d", 123)

	l.WPrint(0, "WPrint Test", 1)
	l.WPrintln(0, "WPrintln Test", 1)
	l.WPrintf("WPrintf Test %d", 123)

	l.EPrint(0, "EPrint Test", 1)
	l.EPrintln(0, "EPrintln Test", 1)
	l.EPrintf("EPrintf Test %d", 123)
}

func Test_debug(t *testing.T) {
	buf := &strings.Builder{}
	loggerOut = buf

	// global logger tests
	{
		SetLevel(LogLevelVerbose)
		VPrint(0, "VPrint Test", 1)
		VPrintln(0, "VPrintln Test", 1)
		VPrintf("VPrintf Test %d", 123)

		IPrint(0, "IPrint Test", 1)
		IPrintln(0, "IPrintln Test", 1)
		IPrintf("IPrintf Test %d", 123)

		WPrint(0, "WPrint Test", 1)
		WPrintln(0, "WPrintln Test", 1)
		WPrintf("WPrintf Test %d", 123)

		EPrint(0, "EPrint Test", 1)
		EPrintln(0, "EPrintln Test", 1)
		EPrintf("EPrintf Test %d", 123)
	}

	logger := NewLogger()
	runLoggerTests(logger)

	taggedLogger := NewLogger("TAG")
	runLoggerTests(taggedLogger)

	taggedLogger = taggedLogger.NewLoggerWithTags("TAG2")
	runLoggerTests(taggedLogger)

	err := Errorf("Err Test")

	logger.EPrintf("%s", err)
	logger.EPrint(err)

	err = ErrorWrapf(err, "Chain test")

	logger.EPrintf("%s", err)
	logger.EPrint(err)

	err = ErrorWrapf(errors.New("Unknown error"), "Unknown error chain test")

	logger.EPrintf("%s", err)
	logger.EPrint(err)

	logger.EPrint(stackTest())
	logger.EPrint(nestedNestedStackTest())

	logger.EPrintf("%s", StackTrace(0))

	r := regexp.MustCompile(output)

	if r.MatchString(buf.String()) {
		t.Log(buf.String())
	} else {
		t.Fatalf("Expected:\n-----\n%s\n-----\nGot:\n-----\n%s-----", output, buf.String())
	}
}

const output = `^\d*:\d*:\d*.\d* \[VERBOSE\] 0VPrint Test1
\d*:\d*:\d*.\d* \[VERBOSE\] 0 VPrintln Test 1
\d*:\d*:\d*.\d* \[VERBOSE\] VPrintf Test 123
\d*:\d*:\d*.\d* \[INFO\]    0IPrint Test1
\d*:\d*:\d*.\d* \[INFO\]    0 IPrintln Test 1
\d*:\d*:\d*.\d* \[INFO\]    IPrintf Test 123
\d*:\d*:\d*.\d* \[WARN\]    0WPrint Test1
\d*:\d*:\d*.\d* \[WARN\]    0 WPrintln Test 1
\d*:\d*:\d*.\d* \[WARN\]    WPrintf Test 123
\d*:\d*:\d*.\d* \[ERROR\]   0EPrint Test1
\d*:\d*:\d*.\d* \[ERROR\]   0 EPrintln Test 1
\d*:\d*:\d*.\d* \[ERROR\]   EPrintf Test 123
\d*:\d*:\d*.\d* \[VERBOSE\] 0VPrint Test1
\d*:\d*:\d*.\d* \[VERBOSE\] 0 VPrintln Test 1
\d*:\d*:\d*.\d* \[VERBOSE\] VPrintf Test 123
\d*:\d*:\d*.\d* \[INFO\]    0IPrint Test1
\d*:\d*:\d*.\d* \[INFO\]    0 IPrintln Test 1
\d*:\d*:\d*.\d* \[INFO\]    IPrintf Test 123
\d*:\d*:\d*.\d* \[WARN\]    0WPrint Test1
\d*:\d*:\d*.\d* \[WARN\]    0 WPrintln Test 1
\d*:\d*:\d*.\d* \[WARN\]    WPrintf Test 123
\d*:\d*:\d*.\d* \[ERROR\]   0EPrint Test1
\d*:\d*:\d*.\d* \[ERROR\]   0 EPrintln Test 1
\d*:\d*:\d*.\d* \[ERROR\]   EPrintf Test 123
\d*:\d*:\d*.\d* \[VERBOSE\] \[TAG\] 0VPrint Test1
\d*:\d*:\d*.\d* \[VERBOSE\] \[TAG\] 0 VPrintln Test 1
\d*:\d*:\d*.\d* \[VERBOSE\] \[TAG\] VPrintf Test 123
\d*:\d*:\d*.\d* \[INFO\]    \[TAG\] 0IPrint Test1
\d*:\d*:\d*.\d* \[INFO\]    \[TAG\] 0 IPrintln Test 1
\d*:\d*:\d*.\d* \[INFO\]    \[TAG\] IPrintf Test 123
\d*:\d*:\d*.\d* \[WARN\]    \[TAG\] 0WPrint Test1
\d*:\d*:\d*.\d* \[WARN\]    \[TAG\] 0 WPrintln Test 1
\d*:\d*:\d*.\d* \[WARN\]    \[TAG\] WPrintf Test 123
\d*:\d*:\d*.\d* \[ERROR\]   \[TAG\] 0EPrint Test1
\d*:\d*:\d*.\d* \[ERROR\]   \[TAG\] 0 EPrintln Test 1
\d*:\d*:\d*.\d* \[ERROR\]   \[TAG\] EPrintf Test 123
\d*:\d*:\d*.\d* \[VERBOSE\] \[TAG\] \[TAG2\] 0VPrint Test1
\d*:\d*:\d*.\d* \[VERBOSE\] \[TAG\] \[TAG2\] 0 VPrintln Test 1
\d*:\d*:\d*.\d* \[VERBOSE\] \[TAG\] \[TAG2\] VPrintf Test 123
\d*:\d*:\d*.\d* \[INFO\]    \[TAG\] \[TAG2\] 0IPrint Test1
\d*:\d*:\d*.\d* \[INFO\]    \[TAG\] \[TAG2\] 0 IPrintln Test 1
\d*:\d*:\d*.\d* \[INFO\]    \[TAG\] \[TAG2\] IPrintf Test 123
\d*:\d*:\d*.\d* \[WARN\]    \[TAG\] \[TAG2\] 0WPrint Test1
\d*:\d*:\d*.\d* \[WARN\]    \[TAG\] \[TAG2\] 0 WPrintln Test 1
\d*:\d*:\d*.\d* \[WARN\]    \[TAG\] \[TAG2\] WPrintf Test 123
\d*:\d*:\d*.\d* \[ERROR\]   \[TAG\] \[TAG2\] 0EPrint Test1
\d*:\d*:\d*.\d* \[ERROR\]   \[TAG\] \[TAG2\] 0 EPrintln Test 1
\d*:\d*:\d*.\d* \[ERROR\]   \[TAG\] \[TAG2\] EPrintf Test 123
\d*:\d*:\d*.\d* \[ERROR\]   Err Test
\d*:\d*:\d*.\d* \[ERROR\]   Err Test
	.*
		.*\S:\d*
\d*:\d*:\d*.\d* \[ERROR\]   Chain test: Err Test
\d*:\d*:\d*.\d* \[ERROR\]   Err Test
	.*
		.*\S:\d*

Chain test
	.*
		.*\S:\d*
\d*:\d*:\d*.\d* \[ERROR\]   Unknown error chain test: Unknown error
\d*:\d*:\d*.\d* \[ERROR\]   Unknown error

Unknown error chain test
	.*
		.*\S:\d*
\d*:\d*:\d*.\d* \[ERROR\]   Stack test
	.*
		.*\S:\d*
	.*
		.*\S:\d*
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
\d*:\d*:\d*.\d* \[ERROR\]   .*
	.*\S:\d*
$`
