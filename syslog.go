package log

import (
	"log/syslog"
)

// Write log messages to syslog.
type SyslogTarget struct {
	writer *syslog.Writer
}

// Create a local syslog writer.
func NewSyslogTarget() (*SyslogTarget, error) {
	if writer, err := syslog.New(syslog.LOG_INFO, prog()); err == nil {
		return &SyslogTarget{writer}, nil
	} else {
		return nil, err
	}
}

// Create a remote syslog writer.
func NewRemoteTarget(network, raddr string) (*SyslogTarget, error) {
	if writer, err := syslog.Dial(network, raddr, syslog.LOG_INFO, prog()); err == nil {
		return &SyslogTarget{writer}, nil
	} else {
		return nil, err
	}
}

// Print a log message to the writer.
func (s *SyslogTarget) Write(level Level, message string) {
	if level == LevelDebug {
		s.writer.Debug(message)
	} else if level == LevelError {
		s.writer.Err(message)
	} else {
		s.writer.Info(message)
	}
}

// Close the writer.
func (s *SyslogTarget) Close() error {
	return s.writer.Close()
}
