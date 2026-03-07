package text

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Config holds text-specific output configuration.
type Config struct {
	RecordSeparator string // Separator between records (default: "\n")
	FieldSeparator  string // Separator between fields (default: " ")
	Template        string // Optional Go template for formatting
	IncludeKeys     bool   // Include field names in output
}

// DefaultConfig returns text emitter defaults.
func DefaultConfig() Config {
	return Config{
		RecordSeparator: "\n",
		FieldSeparator:  " ",
		Template:        "",
		IncludeKeys:     false,
	}
}

// Writer wraps bufio.Writer for text output.
type Writer struct {
	w      *bufio.Writer
	config Config
	first  bool
}

// NewWriter creates a new text writer with the given config.
func NewWriter(w io.Writer, config Config) *Writer {
	return &Writer{
		w:      bufio.NewWriter(w),
		config: config,
		first:  true,
	}
}

// WriteRecord writes a single text record from the canonical payload.
func (w *Writer) WriteRecord(payload map[string]interface{}) error {
	if !w.first {
		if _, err := w.w.WriteString(w.config.RecordSeparator); err != nil {
			return err
		}
	}
	w.first = false

	if w.config.Template != "" {
		// Use template if provided
		return w.writeTemplated(payload)
	}

	// Default: write fields separated by field separator
	var fields []string
	for key, value := range payload {
		var field string
		if w.config.IncludeKeys {
			field = fmt.Sprintf("%s=%v", key, formatValue(value))
		} else {
			field = fmt.Sprintf("%v", formatValue(value))
		}
		fields = append(fields, field)
	}

	_, err := w.w.WriteString(strings.Join(fields, w.config.FieldSeparator))
	return err
}

// Flush ensures all data is written.
func (w *Writer) Flush() error {
	return w.w.Flush()
}

// writeTemplated writes using a template (simplified version).
func (w *Writer) writeTemplated(payload map[string]interface{}) error {
	// For now, use simple key-value substitution
	// A full implementation would use text/template
	output := w.config.Template
	for key, value := range payload {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		output = strings.ReplaceAll(output, placeholder, fmt.Sprintf("%v", formatValue(value)))
	}
	_, err := w.w.WriteString(output)
	return err
}

// formatValue converts interface{} to string representation.
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []interface{}, map[string]interface{}:
		// For complex types, use JSON encoding
		b, _ := json.Marshal(val)
		return string(b)
	default:
		return fmt.Sprintf("%v", val)
	}
}
