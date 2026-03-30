package db_test

import (
	"Monitoring-Pressure/dao/db"
	"testing"

	"gorm.io/gorm"
)

// TestInitDB 测试 InitDB 函数是否能成功初始化数据库
func TestInitDB(t *testing.T) {
	// 调用 InitDB
	db.InitDB()

	if db.DB == nil {
		t.Fatal("DB 全局实例未初始化")
	}

	// 测试能否 ping 数据库
	sqlDB, err := db.DB.DB()
	if err != nil {
		t.Fatalf("获取 sql.DB 失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("数据库无法连接: %v", err)
	}

	// 测试类型是否正确
	var _ *gorm.DB = db.DB
}
