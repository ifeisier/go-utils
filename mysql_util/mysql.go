package mysql_util

import (
	"database/sql"
	"errors"
	"reflect"
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
// db 是 mysql 的连接实例,
// 通过这个二级指针返回新的 mysql 连接实例, 还会关闭旧的连接实例.
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

	if !reflect.DeepEqual(*db, &sql.DB{}) {
		defer func(oldDB *sql.DB) {
			err = oldDB.Close()
		}(*db)
	}
	*db = newDB

	return
}
