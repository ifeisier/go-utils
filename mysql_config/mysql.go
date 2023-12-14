package mysql_config

import (
	"database/sql"
	"time"
)

type Config struct {
	Dsn             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxIdleTime int64  `yaml:"connMaxIdleTime"`
}

func CreateDB(config *Config) (db *sql.DB, err error) {
	dsn := config.Dsn
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime))
	return
}
