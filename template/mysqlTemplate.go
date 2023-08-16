package template

var MysqlTemplate = `package sysinit

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func MysqlInit(setting *AppConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.MysqlC.User,
		setting.MysqlC.Password,
		setting.MysqlC.Host,
		setting.MysqlC.Port,
		setting.MysqlC.Dbname,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	if setting.Mode == "dev" {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: newLogger,
		})
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: newLogger,
		})
	}

	if err != nil {
		if setting.Mode == "dev" {
			fmt.Println(err)
		}
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		if setting.Mode == "dev" {
			fmt.Println(err)
		}
		zap.L().Error("db.DB() failed", zap.Error(err))
	}

	// SetMaxIdle 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(setting.MysqlC.MaxIdle)

	// SetMaxOpen 设置打开数据库的最大数量
	sqlDB.SetMaxOpenConns(setting.MysqlC.MaxOpen)

	return
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("db.DB() failed", zap.Error(err))
	}
	_ = sqlDB.Close()
}
`
