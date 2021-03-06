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
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	LogLevelGlobal = iota
	LogLevelVerbose
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger struct {
	level *uint32
	tags  string
}

var (
	logger    = Logger{level: new(uint32)}
	loggerMtx sync.Mutex
	loggerOut io.StringWriter = os.Stderr
)

func NewLogger(tags ...string) *Logger {
	newLogger := Logger{level: new(uint32)}
	for _, t := range tags {
		newLogger.tags += "[" + t + "] "
	}

	return &newLogger
}

/*
	SetLevel sets the global log level.
*/
func SetLevel(level uint32) {
	if level < LogLevelVerbose || level > LogLevelError {
		panic("Log level out of range")
	}
	atomic.StoreUint32(logger.level, level)
}

/*
	WillLog reports whether logging at the given level would have an effect,
	based on the global log level.
*/
func WillLog(level uint32) bool {
	return atomic.LoadUint32(logger.level) <= level
}

func VPrint(args ...interface{}) {
	logger.VPrint(args...)
}

func IPrint(args ...interface{}) {
	logger.IPrint(args...)
}

func WPrint(args ...interface{}) {
	logger.WPrint(args...)
}

func EPrint(args ...interface{}) {
	logger.EPrint(args...)
}

func VPrintln(args ...interface{}) {
	logger.VPrintln(args...)
}

func IPrintln(args ...interface{}) {
	logger.IPrintln(args...)
}

func WPrintln(args ...interface{}) {
	logger.WPrintln(args...)
}

func EPrintln(args ...interface{}) {
	logger.EPrintln(args...)
}

func VPrintf(format string, args ...interface{}) {
	logger.VPrintf(format, args...)
}

func IPrintf(format string, args ...interface{}) {
	logger.IPrintf(format, args...)
}

func WPrintf(format string, args ...interface{}) {
	logger.WPrintf(format, args...)
}

func EPrintf(format string, args ...interface{}) {
	logger.EPrintf(format, args...)
}

/*
	SetLevel sets the current logger's level, if the logger's level is
	LogLevelGlobal, it will use the global log level.
*/
func (l *Logger) SetLevel(level uint32) {
	if level > LogLevelError {
		panic("Log level out of range")
	}
	atomic.StoreUint32(l.level, level)
}

/*
	WillLog reports whether logging at the given level would have an effect. If
	the logger's level is LogLevelGlobal, it will check against the global log level.
*/
func (l *Logger) WillLog(level uint32) bool {
	logLevel := atomic.LoadUint32(l.level)
	if logLevel == LogLevelGlobal {
		return WillLog(level)
	}
	return logLevel <= level
}

/*
	NewLoggerWithTags creates a new logger with tags appended to the
	current logger's tag list.
*/
func (l *Logger) NewLoggerWithTags(tags ...string) *Logger {
	newLogger := Logger{
		level: new(uint32),
		tags:  l.tags,
	}

	for _, t := range tags {
		newLogger.tags += "[" + t + "] "
	}

	return &newLogger
}

func (l *Logger) messageHeader(level string) string {
	t := time.Now()
	hour, min, sec := t.Clock()

	return fmt.Sprintf("%02d:%02d:%02d.%06d %-10s%s", hour, min, sec, t.Nanosecond()/int(time.Microsecond), "["+level+"]", l.tags)
}

func (l *Logger) VPrint(args ...interface{}) {
	if !l.WillLog(LogLevelVerbose) {
		return
	}

	msg := fmt.Sprintf("%s%s\n", l.messageHeader("VERBOSE"), fmt.Sprint(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) IPrint(args ...interface{}) {
	if !l.WillLog(LogLevelInfo) {
		return
	}

	msg := fmt.Sprintf("%s%s\n", l.messageHeader("INFO"), fmt.Sprint(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) WPrint(args ...interface{}) {
	if !l.WillLog(LogLevelWarn) {
		return
	}

	msg := fmt.Sprintf("%s%s\n", l.messageHeader("WARN"), fmt.Sprint(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) EPrint(args ...interface{}) {
	msg := fmt.Sprintf("%s%s\n", l.messageHeader("ERROR"), fmt.Sprint(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) VPrintln(args ...interface{}) {
	if !l.WillLog(LogLevelVerbose) {
		return
	}

	msg := fmt.Sprintf("%s%s", l.messageHeader("VERBOSE"), fmt.Sprintln(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) IPrintln(args ...interface{}) {
	if !l.WillLog(LogLevelInfo) {
		return
	}

	msg := fmt.Sprintf("%s%s", l.messageHeader("INFO"), fmt.Sprintln(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) WPrintln(args ...interface{}) {
	if !l.WillLog(LogLevelWarn) {
		return
	}

	msg := fmt.Sprintf("%s%s", l.messageHeader("WARN"), fmt.Sprintln(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) EPrintln(args ...interface{}) {
	msg := fmt.Sprintf("%s%s", l.messageHeader("ERROR"), fmt.Sprintln(args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) VPrintf(format string, args ...interface{}) {
	if !l.WillLog(LogLevelVerbose) {
		return
	}

	msg := fmt.Sprintf("%s%s\n", l.messageHeader("VERBOSE"), fmt.Sprintf(format, args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) IPrintf(format string, args ...interface{}) {
	if !l.WillLog(LogLevelInfo) {
		return
	}

	msg := fmt.Sprintf("%s%s\n", l.messageHeader("INFO"), fmt.Sprintf(format, args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) WPrintf(format string, args ...interface{}) {
	if !l.WillLog(LogLevelWarn) {
		return
	}

	msg := fmt.Sprintf("%s%s\n", l.messageHeader("WARN"), fmt.Sprintf(format, args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}

func (l *Logger) EPrintf(format string, args ...interface{}) {
	msg := fmt.Sprintf("%s%s\n", l.messageHeader("ERROR"), fmt.Sprintf(format, args...))

	loggerMtx.Lock()
	defer loggerMtx.Unlock()

	_, _ = loggerOut.WriteString(msg)
}
