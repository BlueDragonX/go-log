package log

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
)

const (
	TargetStderr = "stderr"
	TargetStdout = "stdout"
	TargetSyslog = "syslog"
)

type Option func(logger *Logger) error

// Convert a string to a level value.
func NewLevel(level string) Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LevelDebug
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// Targets are where the logger sends filtered log messages. They can be created by a TargetOpt.
type Target interface {
	Write(level Level, message string)
	Close() error
}

func NewTarget(uri string) (Target, error) {
	var target Target
	var err error
	if uri == TargetStderr {
		target = NewFileTarget(os.Stderr)
	} else if uri == TargetStdout {
		target = NewFileTarget(os.Stdout)
	} else if uri == TargetSyslog {
		target, err = NewSyslogTarget()
	} else {
		var network, raddr string
		if network, raddr, err = uriaddr(uri); err == nil {
			if network == "file" {
				target, err = OpenFileTarget(raddr)
			} else {
				target, err = NewRemoteTarget(network, raddr)
			}
		}
	}
	return target, err
}

// Create a level option to configure the logger.
func LevelOpt(level Level) Option {
	return func(logger *Logger) error {
		logger.SetLevel(level)
		return nil
	}
}

// Create a level option to configure the logger.
func NewLevelOpt(level string) Option {
	return LevelOpt(NewLevel(level))
}

// Create a target option to configure the logger.
func TargetOpt(target Target) Option {
	return func(logger *Logger) error {
		logger.SetTarget(target)
		return nil
	}
}

// Create a target option to configure the logger.
func NewTargetOpt(uri string) Option {
	return func(logger *Logger) error {
		target, err := NewTarget(uri)
		if err == nil {
			logger.SetTarget(target)
		}
		return err
	}
}

// Create a syslog option to configure the logger.
var SyslogOpt Option = func(logger *Logger) error {
	target, err := NewSyslogTarget()
	if err == nil {
		logger.target = target
	}
	return err
}

// Create a console target option to configure the logger.
var ConsoleOpt Option = func(logger *Logger) error {
	logger.target = NewFileTarget(os.Stderr)
	return nil
}

type Logger struct {
	target Target
	level  Level
}

// Return a new logger.
func New(options ...Option) (*Logger, error) {
	logger := &Logger{
		target: NewFileTarget(os.Stderr),
		level:  LevelInfo,
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

// Set the target used by the logger.
func (l *Logger) SetTarget(target Target) {
	l.target = target
}

// Set the level of the logger.
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Close the logger.
func (l *Logger) Close() error {
	return l.target.Close()
}

// Log a message at the provided level.
func (l *Logger) Print(level Level, message string) {
	if level >= l.level {
		l.target.Write(level, message)
	}
}

// Log a formatted message at the provided level.
func (l *Logger) Printf(level Level, format string, a ...interface{}) {
	l.Print(level, fmt.Sprintf(format, a...))
}

// Log a message at the `error` level and call panic().
func (l *Logger) Panic(message string) {
	l.Print(LevelError, message)
	panic(errors.New(message))
}

// Log a formatted message at the `error` level and call panic().
func (l *Logger) Panicf(format string, a ...interface{}) {
	l.Panic(fmt.Sprintf(format, a...))
}

// Log a message at the `error` level and call os.Exit(1).
func (l *Logger) Fatal(message string) {
	l.Print(LevelError, message)
	os.Exit(1)
}

// Log a formatted message at the `error` level and call os.Exit(1).
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

// Log a message at the `debug` level.
func (l *Logger) Debug(message string) {
	l.Print(LevelDebug, message)
}

// Log a formatted message at the `debug` level.
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Printf(LevelDebug, format, a...)
}

// Log a message at the `info` level.
func (l *Logger) Info(message string) {
	l.Print(LevelInfo, message)
}

// Log a formatted message at the `info` level.
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Printf(LevelInfo, format, a...)
}

// Log a message at the `error` level.
func (l *Logger) Error(message string) {
	l.Print(LevelError, message)
}

// Log a formatted message at the `error` level.
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Printf(LevelError, format, a...)
}
