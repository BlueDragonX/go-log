package log

import (
	"testing"
)

type MockWriter struct {
	Level   int
	Message string
}

func NewMockWriter() *MockWriter {
	return &MockWriter{-1, ""}
}

func (w *MockWriter) Write(level int, message string) {
	w.Level = level
	w.Message = message
}

func (w *MockWriter) Close() error {
	return nil
}

func (w *MockWriter) Clear() {
	w.Level = -1
	w.Message = ""
}

func (w *MockWriter) NoWrite() bool {
	nowrite := w.Level == -1 && w.Message == ""
	w.Clear()
	return nowrite
}

func (w *MockWriter) Written(level int, message string) bool {
	written := w.Level == level && w.Message == message
	w.Clear()
	return written
}

func TestPrint(t *testing.T) {
	writer := NewMockWriter()
	logger, _ := New()
	logger.SetWriter(writer)

	logger.Print(LEVEL_DEBUG, "test")
	if !writer.NoWrite() {
		t.Error("debug message written")
	}

	logger.Print(LEVEL_INFO, "test")
	if !writer.Written(LEVEL_INFO, "test") {
		t.Error("info message not written")
	}

	logger.Print(LEVEL_ERROR, "test")
	if !writer.Written(LEVEL_ERROR, "test") {
		t.Error("error message not written")
	}
}

func TestLevel(t *testing.T) {
	writer := NewMockWriter()
	logger, _ := New()
	logger.SetWriter(writer)

	logger.SetLevel(LEVEL_DEBUG)
	logger.Print(LEVEL_DEBUG, "test")
	if !writer.Written(LEVEL_DEBUG, "test") {
		t.Error("debug message not written")
	}

	logger.SetLevel(LEVEL_ERROR)
	logger.Print(LEVEL_DEBUG, "test")
	if !writer.NoWrite() {
		t.Error("debug message written")
	}
}
