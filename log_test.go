package log

import (
	"testing"
)

type MockTarget struct {
	Level   Level
	Message string
}

func NewMockTarget() *MockTarget {
	return &MockTarget{-1, ""}
}

func (w *MockTarget) Write(level Level, message string) {
	w.Level = level
	w.Message = message
}

func (w *MockTarget) Close() error {
	return nil
}

func (w *MockTarget) Clear() {
	w.Level = -1
	w.Message = ""
}

func (w *MockTarget) NoWrite() bool {
	nowrite := w.Level == -1 && w.Message == ""
	w.Clear()
	return nowrite
}

func (w *MockTarget) Written(level Level, message string) bool {
	written := w.Level == level && w.Message == message
	w.Clear()
	return written
}

func TestPrint(t *testing.T) {
	target := NewMockTarget()
	logger, _ := New()
	logger.SetTarget(target)

	logger.Print(LevelDebug, "test")
	if !target.NoWrite() {
		t.Error("debug message written")
	}

	logger.Print(LevelInfo, "test")
	if !target.Written(LevelInfo, "test") {
		t.Error("info message not written")
	}

	logger.Print(LevelError, "test")
	if !target.Written(LevelError, "test") {
		t.Error("error message not written")
	}
}

func TestLevel(t *testing.T) {
	target := NewMockTarget()
	logger, _ := New()
	logger.SetTarget(target)

	logger.SetLevel(LevelDebug)
	logger.Print(LevelDebug, "test")
	if !target.Written(LevelDebug, "test") {
		t.Error("debug message not written")
	}

	logger.SetLevel(LevelError)
	logger.Print(LevelDebug, "test")
	if !target.NoWrite() {
		t.Error("debug message written")
	}
}
