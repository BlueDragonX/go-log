package log

// Writer is an io.Writer which logs to Logger at the given Level.
type Writer struct {
	Level  Level
	Logger *Logger
}

// NewWriter creates a new Writer.
func NewWriter(level Level, logger *Logger) *Writer {
	return &Writer{level, logger}
}

// Write msg to the logger.
func (w *Writer) Write(msg []byte) (int, error) {
	w.Logger.Printf(w.Level, "%s", msg)
	return len(msg), nil
}
