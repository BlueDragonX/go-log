package log

import (
	"io"
	"testing"
)

func TestWriter(t *testing.T) {
	target := NewMockTarget()
	logger, _ := New()
	logger.SetTarget(target)

	var wrt io.Writer = NewWriter(LevelDebug, logger)
	msg := "test"

	wrt.Write([]byte(msg))
	if !target.NoWrite() {
		t.Error("debug message written")
	}

	logger.SetLevel(LevelDebug)
	n, err := wrt.Write([]byte(msg))
	if err != nil {
		t.Error(err)
	} else if n != len(msg) {
		t.Error("invalid message length")
	} else if !target.Written(LevelDebug, msg) {
		t.Error("debug message not written")
	}
}
