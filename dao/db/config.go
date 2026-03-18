package db

import (
	"Monitoring-Pressure/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
) // 全局DB实例

// InitDB 初始化数据库连接
func InitDB() {
	// 数据库连接信息（建议通过环境变量/配置文件读取，这里为示例硬编码）
	dsn := ""

	// 自定义日志配置（开发环境显示SQL，生产环境可关闭）
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 日志级别：Silent/Error/Warn/Info
			Colorful:      true,        // 彩色打印
		},
	)

	// 连接数据库并配置连接池
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
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	err = models.AutoMigrate(db)
	if err != nil {
		panic(err)
	}

	DB = db
	log.Println("数据库连接成功")
}
