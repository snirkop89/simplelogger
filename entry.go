package simplelogger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Entry struct {
	minLevel   LogLevel
	format     int
	out        io.Writer
	Level      string            `json:"level"`
	Time       string            `json:"time"`
	Message    string            `json:"message"`
	Properties map[string]string `json:"properties,omitempty"`
	Trace      string            `json:"trace,omitempty"`
}

func newEntry(out io.Writer, minLevel LogLevel, format int) *Entry {
	return &Entry{
		minLevel:   minLevel,
		out:        out,
		format:     format,
		Time:       time.Now().UTC().Format(time.RFC3339),
		Properties: make(map[string]string),
	}
}

func (e *Entry) Info(msg string) {
	e.print(LevelInfo, msg)
}

func (e *Entry) Infof(msg string, args ...any) {
	e.print(LevelInfo, msg, args...)
}

func (e *Entry) Warn(msg string) {
	e.print(LevelWarn, msg)
}

func (e *Entry) Warnf(msg string, args ...any) {
	e.print(LevelWarn, msg, args...)
}

func (e *Entry) Error(msg string) {
	e.print(LevelError, msg)
}

func (e *Entry) Errorf(msg string, args ...any) {
	e.print(LevelError, msg, args...)
}

func (e *Entry) print(level LogLevel, msg string, args ...any) {
	if level < e.minLevel {
		return
	}

	e.Level = level.String()
	e.Message = fmt.Sprintf(msg, args...)

	if e.format == FormatHuman {
		cw := &consoleWriter{entry: e}
		cw.print()
		return
	}

	line, err := json.Marshal(e)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal message: " + err.Error())
	}

	e.Write(line)
}

func (e *Entry) Write(data []byte) {
	fmt.Fprintln(e.out, string(data))
}

type consoleWriter struct {
	entry *Entry
}

func (cw *consoleWriter) print() {
	msg := fmt.Sprintf("%s\t%s\t%s\t",
		strings.ToUpper(cw.entry.Level),
		cw.entry.Time, cw.entry.Message)

	var builder strings.Builder
	for k, v := range cw.entry.Properties {
		builder.WriteString(fmt.Sprintf("%s=%v ", k, v))
	}
	msg += builder.String() + "\n"
	cw.Write(msg)
}

func (cw *consoleWriter) Write(msg string) (int, error) {
	return fmt.Fprint(os.Stdout, msg)
}
