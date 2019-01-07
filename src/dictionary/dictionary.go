package dictionary

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Dictionary struct {
	Source      string   `yaml:"source"`
	Table       string   `yaml:"table"`
	Columns     []string `yaml:"columns"`
	Where       string   `yaml:"where"`
	UpdateField string   `yaml:"update_field"`
}

// Handler make an HTTP handler
func Handler(dict *Dictionary) (http.HandlerFunc, error) {
	const updatedAtLayout = "2006-01-02 15:04:05"
	var (
		query           = fmt.Sprintf("SELECT %[1]s FROM %[2]s", strings.Join(dict.Columns, ", "), dict.Table)
		dictWhere       = dict.Where
		dictUpdateField = dict.UpdateField
		connect, found  = sources[dict.Source]
	)

	if !found {
		return nil, fmt.Errorf("source '%s' not found", dict.Source)
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		var where []string
		if len(dictWhere) != 0 {
			where = append(where, dictWhere)
		}
		if updatedAt := req.URL.Query().Get(dictUpdateField); len(updatedAt) != 0 {
			time, err := time.Parse(updatedAtLayout, updatedAt)
			if err != nil {
				log.Printf("invalid time format: %v", err)
				return
			}
			where = append(where, fmt.Sprintf(dictUpdateField+" > '%s'", time.Format(updatedAtLayout)))
		}
		querySQL := query
		if len(where) != 0 {
			querySQL = querySQL + " WHERE " + strings.Join(where, " AND ")
		}
		rows, err := connect.Query(querySQL)
		if err != nil {
			log.Println(err)
			return
		}
		defer rows.Close()
		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			log.Println(err)
			return
		}
		var (
			columns = make([]string, 0, len(columnTypes))
			values  = make([]interface{}, 0, len(columnTypes))
		)

		for _, columnType := range columnTypes {
			columns = append(columns, columnType.Name())
			values = append(values, reflect.New(columnType.ScanType()).Interface())
		}

		output := FormatFactory(req.URL.Query().Get("format"), columns, rw)

		for rows.Next() {
			if err := rows.Scan(values...); err != nil {
				log.Println(err)
				return
			}
			output.Write(values)
		}
		output.Flush()
	}, nil
}
