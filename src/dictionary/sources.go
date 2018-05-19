package dictionary

import (
	"database/sql"
	"fmt"
	"time"
)

var sources = make(map[string]*sql.DB)

// RegisterSources ...
func RegisterSources(list []Source) (err error) {
	for _, source := range list {
		switch driver, name := source.Source, source.Name; driver {
		case "mysql", "postgres":
			if _, duplicate := sources[name]; duplicate {
				return fmt.Errorf("source '%s' already exists", name)
			}
			if sources[name], err = sql.Open(driver, source.DSN); err != nil {
				return fmt.Errorf("could not open source '%s': %v", name, err)
			}
			sources[name].SetMaxIdleConns(source.MaxIdleConns)
			sources[name].SetMaxOpenConns(source.MaxOpenConns)
			sources[name].SetConnMaxLifetime(source.ConnMaxLifetime)
		default:
			return fmt.Errorf("unknow source '%s'", name)
		}
	}
	return nil
}

// Source is a credentials for the source database
type Source struct {
	Source          string        `yaml:"source"`
	DSN             string        `yaml:"dsn"`
	Name            string        `yaml:"name"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}
