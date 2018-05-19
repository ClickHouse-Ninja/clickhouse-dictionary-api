package dictionary_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/ClickHouse-Ninja/clickhouse-dictionary-api/src/dictionary"
	"github.com/stretchr/testify/assert"
)

func TestTabSeparated(t *testing.T) {
	var (
		buffer  bytes.Buffer
		format  = dictionary.FormatFactory("TabSeparated", []string{}, &buffer)
		date, _ = time.Parse("2006-01-02", "2018-05-18")
	)

	if err := format.Write([]interface{}{"a", 1, 1.5, date}); assert.NoError(t, err) {
		format.Flush()
		assert.Equal(t, "a\t1\t1.5\t2018-05-18\n", buffer.String())
	}
}

func TestJSONEachRow(t *testing.T) {
	var (
		buffer  bytes.Buffer
		format  = dictionary.FormatFactory("JSONEachRow", []string{"string", "int", "float", "date"}, &buffer)
		date, _ = time.Parse("2006-01-02", "2018-05-18")
	)
	if err := format.Write([]interface{}{"a", 1, 1.5, date}); assert.NoError(t, err) {
		assert.Equal(t, "{\"date\":\"2018-05-18\",\"float\":1.5,\"int\":1,\"string\":\"a\"}\n", buffer.String())
	}
}
