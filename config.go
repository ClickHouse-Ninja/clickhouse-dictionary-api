package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var config struct {
	Sources []struct {
		Source string `yaml:"source"`
		DSN    string `yaml:"dsn"`
		Name   string `yaml:"name"`
	} `yaml:"sources"`
	Dictionaries map[string]dictionary `yaml:"dictionaries"`
}

type dictionary struct {
	Source  string   `yaml:"source"`
	Table   string   `yaml:"table"`
	Columns []string `yaml:"columns"`
	Where   string   `yaml:"where"`
}

var (
	listenAddr = flag.String("addr", ":8080", "")
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
