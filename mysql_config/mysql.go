package mysql_config

import (
	"database/sql"
	"errors"
	"time"
)

type Config struct {
	Dsn             string `yaml:"dsn" json:"dsn"`
	MaxOpenConn     int    `yaml:"maxOpenConn" json:"maxOpenConn"`
	ConnMaxIdleTime int64  `yaml:"connMaxIdleTime" json:"connMaxIdleTime"`
}

// DefaultMysqlDB 获取空的 MySQL 连接
func DefaultMysqlDB() **sql.DB {
	ptrDB := &sql.DB{}
	return &ptrDB
}

// CreateDB 创建 MySQL 的数据库连接
//
// config 是 MySQL 连接的配置
//
// db 是一个二级指针, 通过这个二级指针返回新的 MySQL 连接实例,
// 如果给定的 db 是已经连接的, 当创建新的连接实例后, 会关闭旧连接实例.
func CreateDB(config *Config, db **sql.DB) (err error) {
	newDB, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		return err
	}
	newDB.SetMaxOpenConns(config.MaxOpenConn)
	newDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime))

	if db == nil {
		return errors.New("db 不能是 nil")
	}

	if *db != (&sql.DB{}) {
		defer func(oldDB *sql.DB) {
			err = oldDB.Close()
		}(*db)
	}
	*db = newDB

	return
}
