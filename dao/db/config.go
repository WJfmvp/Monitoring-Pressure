package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"Monitoring-Pressure/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB // 全局 DB 实例

// InitDB 初始化数据库连接
func InitDB() {
	// 显式加载 .env 文件，确保从文件读取配置
	err := godotenv.Load("F:/Monitoring-Pressure/.env")
	if err != nil {
		panic(fmt.Sprintf("未找到 .env 文件或加载失败: %v", err)) // 强制抛出错误
	}

	// 获取 DB_DNS
	dsn := os.Getenv("DB_DNS")
	if dsn == "" {
		panic("环境变量 DB_DNS 未设置或加载失败")
	}

	// 打印 DSN 用于调试
	fmt.Println("连接数据库 DSN:", dsn)

	// 自定义日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接池失败: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	if err := models.Migrate(db); err != nil {
		panic(err)
	}

	DB = db
	log.Println("数据库连接成功")
}
