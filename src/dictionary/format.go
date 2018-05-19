package dictionary

import (
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// FormatFactory ...
func FormatFactory(format string, columns []string, w io.Writer) Format {
	switch format {
	case "TabSeparated":
		writer := csv.NewWriter(w)
		writer.Comma = '\t'
		return &TabSeparated{
			writer: writer,
		}
	default:
		return &JSONEachRow{
			columns: columns,
			encoder: json.NewEncoder(w),
		}
	}
}

// Format output
type Format interface {
	Write(values []interface{}) error
	Flush()
}

// TabSeparated output
type TabSeparated struct {
	writer *csv.Writer
}

// Write row to the output
func (tab *TabSeparated) Write(values []interface{}) error {
	record := make([]string, 0, len(values))
	for _, value := range values {
		record = append(record, fmt.Sprint(convert(value)))
	}
	return tab.writer.Write(record)
}

// Flush writes
func (tab *TabSeparated) Flush() {
	tab.writer.Flush()
}

// JSONEachRow output
type JSONEachRow struct {
	columns []string
	encoder *json.Encoder
}

// Write row to the output
func (json *JSONEachRow) Write(values []interface{}) error {
	if len(json.columns) != len(values) {
		return fmt.Errorf("unexpected number of columns: expected %d, got %d", len(json.columns), len(values))
	}
	row := make(map[string]interface{}, len(json.columns))
	for i, value := range values {
		row[json.columns[i]] = convert(value)
	}
	return json.encoder.Encode(row)
}

// Flush ...
func (JSONEachRow) Flush() {}

func convert(value interface{}) interface{} {
	switch value := value.(type) {
	case *sql.RawBytes:
		return string(*value)
	case driver.Valuer:
		v, _ := value.Value()
		return v
	case time.Time:
		return value.Format("2006-01-02")
	case *time.Time:
		return value.Format("2006-01-02")
	default:
		return value
	}
}
