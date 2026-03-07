package json

import (
	"io"
)

// Config holds JSON-specific output configuration.
type Config struct {
	// Future: Pretty print, indentation, etc.
}

// DefaultConfig returns JSON emitter defaults.
func DefaultConfig() Config {
	return Config{}
}

// Writer wraps an io.Writer for JSON line-delimited output.
type Writer struct {
	w io.Writer
}

// NewWriter creates a new JSON writer.
func NewWriter(w io.Writer, config Config) *Writer {
	return &Writer{w: w}
}

// WriteRecord writes a single JSON record (already serialized as []byte).
// For Phase 1, we simply write the JSON bytes from omniparser with a newline.
func (w *Writer) WriteRecord(jsonBytes []byte) error {
	if _, err := w.w.Write(jsonBytes); err != nil {
		return err
	}
	if _, err := w.w.Write([]byte{'\n'}); err != nil {
		return err
	}
	return nil
}
