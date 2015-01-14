package log

import (
	"log/syslog"
)

// Write log messages to syslog.
type SyslogWriter struct {
	writer *syslog.Writer
}

// Create a local syslog writer.
func NewSyslogWriter() (*SyslogWriter, error) {
	if writer, err := syslog.New(syslog.LOG_INFO, prog()); err == nil {
		return &SyslogWriter{writer}, nil
	} else {
		return nil, err
	}
}

// Create a remote syslog writer.
func NewRemoteWriter(network, raddr string) (*SyslogWriter, error) {
	if writer, err := syslog.Dial(network, raddr, syslog.LOG_INFO, prog()); err == nil {
		return &SyslogWriter{writer}, nil
	} else {
		return nil, err
	}
}

// Print a log message to the writer.
func (s *SyslogWriter) Write(level int, message string) {
	if level == LEVEL_DEBUG {
		s.writer.Debug(message)
	} else if level == LEVEL_ERROR {
		s.writer.Err(message)
	} else {
		s.writer.Info(message)
	}
}

// Close the writer.
func (s *SyslogWriter) Close() error {
	return s.writer.Close()
}
