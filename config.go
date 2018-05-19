package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/ClickHouse-Ninja/clickhouse-dictionary-api/src/dictionary"
	"gopkg.in/yaml.v2"
)

var config struct {
	Address      string                           `yaml:"address"`
	Sources      []dictionary.Source              `yaml:"sources"`
	Dictionaries map[string]dictionary.Dictionary `yaml:"dictionaries"`
}

var (
	configPath = flag.String("c", "config.yaml", "")
)

func init() {
	flag.Parse()
	{
		data, err := ioutil.ReadFile(*configPath)
		if err != nil {
			log.Fatalf("could not open config file: %v", err)
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Fatalf("could not read config file: %v", err)
		}
	}
}
