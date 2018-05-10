package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var sources = make(map[string]*sql.DB)

func main() {
	for _, source := range config.Sources {
		if _, exists := sources[source.Name]; exists {
			log.Fatalf("source '%s' already exists", source.Name)
		}
		var err error
		switch source.Source {
		case "mysql", "postgres":
			if sources[source.Name], err = sql.Open(source.Source, source.DSN); err != nil {
				log.Fatalf("could not open '%s': %v", source.Source, err)
			}
			sources[source.Name].SetMaxOpenConns(10)
			sources[source.Name].SetMaxIdleConns(2)
			sources[source.Name].SetConnMaxLifetime(5 * time.Minute)
		default:
			log.Fatalf("unknow source '%s'", source.Source)
		}
	}

	for name, dictionary := range config.Dictionaries {
		if handler, err := makeHandler(&dictionary); err == nil {
			log.Println("add handler: /dictionary/" + name)
			http.HandleFunc("/dictionary/"+name, handler)
		} else {
			log.Fatalf("could not add dictionary handler: %v", err)
		}
	}

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("could not listen address '%s': %v", "ss", err)
	}
}

func makeHandler(dictionary *dictionary) (http.HandlerFunc, error) {
	var where string
	if len(dictionary.Where) != 0 {
		where = "WHERE " + dictionary.Where
	}
	sqlQuery := fmt.Sprintf(`SELECT %[1]s FROM %[2]s %[3]s`,
		strings.Join(dictionary.Columns, ", "),
		dictionary.Table,
		where,
	)
	connect, found := sources[dictionary.Source]
	if !found {
		return nil, fmt.Errorf("source '%s' not found", dictionary.Source)
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		rows, err := connect.Query(sqlQuery)
		if err != nil {
			log.Printf("could not execute query: %v", err)
			return
		}

		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			log.Printf("could not read columns types: %v", err)
			return
		}

		columnNames, err := rows.Columns()
		if err != nil {
			log.Printf("could not read columns: %v", err)
			return
		}

		if len(columnTypes) != len(columnNames) {
			log.Fatalf("unexpected number of columns: expected %d, got %d", len(columnNames), len(columnTypes))
			return
		}

		values := make([]interface{}, 0, len(columnTypes))
		for _, columnType := range columnTypes {
			values = append(values, reflect.New(columnType.ScanType()).Interface())
		}

		encoder := json.NewEncoder(rw)
		for rows.Next() {
			if err := rows.Scan(values...); err != nil {
				log.Printf("could not scan values: %v", err)
				return
			}
			row := map[string]interface{}{}
			for i, column := range columnNames {
				switch v := values[i].(type) {
				case *sql.RawBytes:
					row[column] = string(*v)
				case driver.Valuer:
					value, _ := v.Value()
					row[column] = value
				default:
					row[column] = values[i]
				}
			}
			encoder.Encode(row)
		}
	}, nil
}
