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
	"log"
	"regexp"
	"strings"
	"testing"
)

func stackTest() error {
	return ErrorNew("Stack test")
}

func passthroughTest() error {
	return stackTest()
}

func nestedStackTest() error {
	return ErrorWrap(ErrorWrap(passthroughTest(), "Nested stack test"), "Double nested test")
}

func nestedNestedStackTest() error {
	return ErrorWrap(nestedStackTest(), "Nested nested stack test")
}

func Test_debug(t *testing.T) {
	buf := &strings.Builder{}
	logOut = log.New(buf, "", 0).Output
	LogSetLevel(LogLevelVerbose)

	LogV("V Test")
	LogI("I Test")
	LogW("W Test")
	LogE("E Test")

	err := ErrorNew("Err Test")

	LogE("%s", err)
	LogErr(err)

	err = ErrorWrap(err, "Chain test")

	LogE("%s", err)
	LogErr(err)

	err = ErrorWrap(errors.New("Unknown error"), "Unknown error chain test")

	LogE("%s", err)
	LogErr(err)

	LogErr(stackTest())
	LogErr(nestedNestedStackTest())

	LogE("%s", StackTrace(0))

	r := regexp.MustCompile(output)

	if !r.MatchString(buf.String()) {
		t.Fatalf("Expected:\n-----\n%s\n-----\nGot:\n-----\n%s-----", output, buf.String())
	}
}

const output = `^\[VERBOSE\] V Test
\[INFO\]    I Test
\[WARN\]    W Test
\[ERROR\]   E Test
\[ERROR\]   Err Test
\[ERROR\]   Err Test
	.*
		.*\S:\d*
\[ERROR\]   Chain test: Err Test
\[ERROR\]   Err Test
	.*
		.*\S:\d*

Chain test
	.*
		.*\S:\d*
\[ERROR\]   Unknown error chain test: Unknown error
\[ERROR\]   Unknown error

Unknown error chain test
	.*
		.*\S:\d*
\[ERROR\]   Stack test
	.*
		.*\S:\d*
	.*
		.*\S:\d*
\[ERROR\]   Stack test
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
\[ERROR\]   .*
	.*\S:\d*
$`
