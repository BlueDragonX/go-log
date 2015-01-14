package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type FileWriter struct {
	prog string
	writer io.Writer
}

// Create a new file writer pointing to the file at the provided path.
func OpenFileWriter(file string) (*FileWriter, error) {
	if writer, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		return &FileWriter{prog: prog(), writer: writer}, nil
	} else {
		return nil, err
	}
}

// Create a new file writer pointing to the provided open file.
func NewFileWriter(writer io.Writer) *FileWriter {
	return &FileWriter{prog: prog(), writer: writer}
}

// Print a log message to the writer.
func (s *FileWriter) Write(level int, message string) {
	now := time.Now().Format("2006/01/02 15:04:05")
	fmt.Fprintf(s.writer, "%s %s: %s\n", now, s.prog, message)
}

// Close the writer.
func (s *FileWriter) Close() error {
	if s.writer == os.Stderr || s.writer == os.Stdout {
		// Don't close stdio writers.
		return nil
	}
	if closer, ok := s.writer.(io.Closer); ok {
		// Close the writer if allowed.
		return closer.Close()
	}
	return nil
}
