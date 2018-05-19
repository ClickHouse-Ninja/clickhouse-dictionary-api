package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ClickHouse-Ninja/clickhouse-dictionary-api/src/dictionary"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var sources = make(map[string]*sql.DB)

func main() {

	if err := dictionary.RegisterSources(config.Sources); err != nil {
		log.Fatalf("an error occurred during registration of the source: %v", err)
	}

	for name, dict := range config.Dictionaries {
		if handler, err := dictionary.Handler(&dict); err == nil {
			log.Println("add handler: /dictionary/" + name)
			http.HandleFunc("/dictionary/"+name, handler)
		} else {
			log.Fatalf("could not add dictionary handler: %v", err)
		}
	}

	if err := http.ListenAndServe(config.Address, nil); err != nil {
		log.Fatalf("could not bind address '%s': %v", config.Address, err)
	}
}
