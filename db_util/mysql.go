package db_util

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"time"
)

type MysqlConfig struct {
	Dsn string `yaml:"dsn" json:"dsn"`
	// 最大连接数
	MaxOpenConn int `yaml:"maxOpenConn" json:"maxOpenConn"`
	// 空闲连接数
	MaxIdleConn int `yaml:"maxIdleConn" json:"maxIdleConn"`
	// ConnMaxIdleTime 连接空间的最长时间,单位是分钟
	ConnMaxIdleTime int64 `yaml:"connMaxIdleTime" json:"connMaxIdleTime"`
}

// DefaultMysqlDB 获取空的 MySQL 连接
func DefaultMysqlDB() **sql.DB {
	ptrDB := &sql.DB{}
	return &ptrDB
}

// CreateMysqlDB 创建 MySQL 的数据库连接
//
// config 是 MySQL 连接的配置
func CreateMysqlDB(config *MysqlConfig) (*sql.DB, error) {
	newDB, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		return nil, err
	}
	newDB.SetMaxOpenConns(config.MaxOpenConn)
	newDB.SetMaxIdleConns(config.MaxIdleConn)
	newDB.SetConnMaxIdleTime(time.Minute * time.Duration(config.ConnMaxIdleTime))
	return newDB, nil
}

// RefreshMysqlDB 刷新 Mysql 数据库连接
// 使用新连接代替旧连接, 并关闭旧连接.
//
// oldDB 需要一个二级指针变量, 会更新这个二级指针的 Mysql 连接
//
// newBD 新的 Mysql 连接
func RefreshMysqlDB(oldDB **sql.DB, newBD *sql.DB) (err error) {
	if oldDB == nil {
		return errors.New("oldDB 不能是 nil")
	}

	if !reflect.DeepEqual(*oldDB, &sql.DB{}) {
		defer func(old *sql.DB) {
			err = old.Close()
		}(*oldDB)
	}
	*oldDB = newBD
	return
}
