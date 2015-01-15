package log

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_ERROR
)

const (
	TARGET_STDERR = "stderr"
	TARGET_STDOUT = "stdout"
	TARGET_SYSLOG = "syslog"
)

type Option func(logger *Logger) error

// Writers are where the logger sends filtered log messages and are created by a target.
type Writer interface {
	Write(level int, message string)
	Close() error
}

// Create a level option to configure the logger.
func Level(level string) Option {
	level = strings.ToLower(level)
	levelInt := LEVEL_INFO
	switch level {
	case "debug":
		levelInt = LEVEL_DEBUG
	case "error":
		levelInt = LEVEL_ERROR
	}

	return func(logger *Logger) error {
		logger.level = levelInt
		return nil
	}
}

// Create a target option to configure the logger.
func Target(uri string) Option {
	var writer Writer
	var err error
	if uri == TARGET_STDERR {
		writer = NewFileWriter(os.Stderr)
	} else if uri == TARGET_STDOUT {
		writer = NewFileWriter(os.Stdout)
	} else if uri == TARGET_SYSLOG {
		writer, err = NewSyslogWriter()
	} else {
		var network, raddr string
		if network, raddr, err = uriaddr(uri); err == nil {
			if network == "file" {
				writer, err = OpenFileWriter(raddr)
			} else {
				writer, err = NewRemoteWriter(network, raddr)
			}
		}
	}

	return func(logger *Logger) error {
		if err == nil {
			logger.writer = writer
		}
		return err
	}
}

// Create a syslog option to configure the logger.
var Syslog Option = func(logger *Logger) error {
	writer, err := NewSyslogWriter()
	if err == nil {
		logger.writer = writer
	}
	return err
}

// Create a console target option to configure the logger.
var Console Option = func(logger *Logger) error {
	logger.writer = NewFileWriter(os.Stderr)
	return nil
}

type Logger struct {
	writer Writer
	level  int
}

// Return a new logger.
func New(options ...Option) (*Logger, error) {
	logger := &Logger{
		writer: NewFileWriter(os.Stderr),
		level:  LEVEL_INFO,
	}

	var err error
	for _, option := range options {
		if err = option(logger); err != nil {
			return nil, err
		}
	}
	return logger, nil
}

// Return a new logger. If an error occurs print it to stderr and call os.Exit(1).
func NewOrExit(options ...Option) *Logger {
	logger, err := New(options...)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}
	return logger
}

// Set the writer used by the logger.
func (l *Logger) SetWriter(writer Writer) {
	l.writer = writer
}

// Set the level of the logger.
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// Close the logger.
func (l *Logger) Close() error {
	return l.writer.Close()
}

// Log a message at the provided level.
func (l *Logger) Print(level int, message string) {
	if level >= l.level {
		l.writer.Write(level, message)
	}
}

// Log a formatted message at the provided level.
func (l *Logger) Printf(level int, format string, a ...interface{}) {
	l.Print(level, fmt.Sprintf(format, a...))
}

// Log a message at the `error` level and call panic().
func (l *Logger) Panic(message string) {
	l.Print(LEVEL_ERROR, message)
	panic(errors.New(message))
}

// Log a formatted message at the `error` level and call panic().
func (l *Logger) Panicf(format string, a ...interface{}) {
	l.Panic(fmt.Sprintf(format, a...))
}

// Log a message at the `error` level and call os.Exit(1).
func (l *Logger) Fatal(message string) {
	l.Print(LEVEL_ERROR, message)
	os.Exit(1)
}

// Log a formatted message at the `error` level and call os.Exit(1).
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

// Log a message at the `debug` level.
func (l *Logger) Debug(message string) {
	l.Print(LEVEL_DEBUG, message)
}

// Log a formatted message at the `debug` level.
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Printf(LEVEL_DEBUG, format, a...)
}

// Log a message at the `info` level.
func (l *Logger) Info(message string) {
	l.Print(LEVEL_INFO, message)
}

// Log a formatted message at the `info` level.
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Printf(LEVEL_INFO, format, a...)
}

// Log a message at the `error` level.
func (l *Logger) Error(message string) {
	l.Print(LEVEL_ERROR, message)
}

// Log a formatted message at the `error` level.
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Printf(LEVEL_ERROR, format, a...)
}
