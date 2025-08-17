package data

import (
	"comment/internal/conf"
	"comment/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(c *conf.Data_Database) (*gorm.DB, error) {
	log.Info(nil, "init gorm db")
	log.Debug(nil, "db config", "db", c)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: c.Source,
	}))
	if err != nil {
		log.Error(nil, "failed to connect database", "error", err)
		return nil, err
	}
	log.Info(nil, "connect database successful.")

	// 设置连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(nil, "failed to get database instance", "error", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(int(c.IdleConns))
	sqlDB.SetMaxOpenConns(int(c.MaxOpenConns))
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime.AsDuration())
	sqlDB.SetConnMaxIdleTime(c.ConnMaxIdleTime.AsDuration())
	log.Info(nil, "set database pool successful.")

	return db, nil
}
