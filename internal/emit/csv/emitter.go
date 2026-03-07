package csv

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/pbagtoltol/omniwriter/internal/schema"
)

// Config holds CSV-specific output configuration.
type Config struct {
	Delimiter rune
	Columns   []schema.CSVColumn
}

// DefaultConfig returns CSV emitter defaults.
func DefaultConfig() Config {
	return Config{
		Delimiter: ',',
	}
}

// Writer wraps csv.Writer with omniwriter-specific configuration.
type Writer struct {
	w       *csv.Writer
	columns []schema.CSVColumn
}

// NewWriter creates a new CSV writer with the given config.
func NewWriter(w io.Writer, config Config) *Writer {
	csvWriter := csv.NewWriter(w)
	if config.Delimiter != 0 {
		csvWriter.Comma = config.Delimiter
	}
	return &Writer{
		w:       csvWriter,
		columns: config.Columns,
	}
}

// WriteRecord writes a single CSV row from the canonical payload.
func (w *Writer) WriteRecord(payload map[string]interface{}) error {
	row := make([]string, 0, len(w.columns))
	for _, col := range w.columns {
		if col.Path != "" {
			row = append(row, stringify(getPath(payload, col.Path)))
		} else {
			row = append(row, col.Const)
		}
	}
	return w.w.Write(row)
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	w.w.Flush()
	return w.w.Error()
}

// getPath retrieves a value from the payload using a dot-separated path.
func getPath(payload map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var cur interface{} = payload
	for _, p := range parts {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur = m[p]
	}
	return cur
}

// stringify converts a value to string representation.
func stringify(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		return formatFloat(t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// formatFloat formats a float64, avoiding scientific notation.
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
