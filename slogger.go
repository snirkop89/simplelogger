// Package slogger is a simple logger that supports log levels
// and two formats: human-redable (tab-delimited values) and stractures (JSON).
//
// It does not have any dependency and can be used for small projects and
// keep the code base / vendor folder small.
package simplelogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type LogLevel int

// LogLevels
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "none"
	}
}

// Format options
const (
	FormatHuman int = iota
	FormatJSON
)

type Opt func(*Logger)

type Logger struct {
	format   int
	minLevel LogLevel
	out      io.Writer
}

func New(format int, minLevel LogLevel, opts ...Opt) *Logger {
	logger := &Logger{
		format:   format,
		minLevel: LogLevel(minLevel),
		out:      os.Stdout,
	}

	for _, o := range opts {
		o(logger)
	}

	return logger
}

func WithWriter(out io.Writer) func(*Logger) {
	return func(l *Logger) {
		l.out = out
	}
}

func (l *Logger) Info(msg string) {
	e := newEntry(l.out, l.minLevel, l.format)
	e.print(LevelInfo, msg)
}

func (l *Logger) Infof(msg string, args ...any) {
	e := newEntry(l.out, l.minLevel, l.format)
	e.print(LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string) {
	e := newEntry(l.out, l.minLevel, l.format)
	e.print(LevelWarn, msg)
}

func (l *Logger) Warnf(msg string, args ...any) {
	e := newEntry(l.out, l.minLevel, l.format)
	e.print(LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string) {
	e := newEntry(l.out, l.minLevel, l.format)
	e.print(LevelError, msg)
}

func (l *Logger) Errorf(msg string, args ...any) {
	e := newEntry(l.out, l.minLevel, l.format)
	e.print(LevelError, msg, args...)
}

func (l *Logger) WithFields(pairs ...string) *Entry {
	e := newEntry(l.out, l.minLevel, l.format)
	if len(pairs) == 0 {
		return e
	}

	// Validate number of pairs
	if len(pairs)%2 != 0 {
		e.Properties["fields_error"] = fmt.Sprintf("invalid number of fields: %d [%s]", len(pairs), strings.Join(pairs, ","))
		return e
	}
	for i := 0; i < len(pairs)-1; i++ {
		e.Properties[pairs[i]] = pairs[i+1]
		i++
	}

	return e
}

// Wrapper for built-in functions, for user-friendly API, and reducing
// the need in importing the "log" package
func Print(msg string) {
	log.Print(msg + "\n")
}

func Println(msg string) {
	log.Println(msg + "\n")
}

func Printf(msg string, args ...any) {
	log.Printf(msg+"\n", args...)
}
