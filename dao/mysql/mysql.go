package mysql

import (
	"fmt"
	"time"

	"github.com/bluebell/models"

	"go.uber.org/zap"

	"github.com/bluebell/settings"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init(cfg *settings.MysqlConfig, mode string) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect to mysql failed", zap.Error(err))
		return
	}
	//开发环境，实时打印SQL
	if mode == "dev" {
		Db = Db.Debug()
	}
	db, err := Db.DB()
	if err != nil {
		zap.L().Error("get Db.DB() failed,err:%v\n", zap.Error(err))
		return
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()

	if err != nil {
		return
	}

	//同步表结构
	err = AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	err = AutoMigrate(&models.Community{})
	if err != nil {
		return err
	}
	err = AutoMigrate(&models.Post{})
	if err != nil {
		return err
	}
	return
}

func AutoMigrate(model interface{}) (err error) {
	err = Db.AutoMigrate(model)
	return err
}
