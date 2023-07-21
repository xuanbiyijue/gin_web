/*
用来管理数据库连接和关闭
 */

package mysql

import (
	"fmt"
	"ginQAP/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)


// db 仅本包可用的db实例
var db *sqlx.DB


// Init 初始化MySQL
func Init(config *config.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DbName,
	)
	//zap.L().Info("dsn", zap.String("dsn", dsn))
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect MySQL failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	return
}

// Close 关闭数据库连接
func Close() {
	_ = db.Close()
}