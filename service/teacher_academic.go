package service

import (
	"fmt"
	"time"

	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/models/users"
	"gorm.io/gorm"
)

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// GetAcademicRecordsAdvanced 高级查询（分页 + 筛选 + 时间范围）
func GetAcademicRecordsAdvanced(
	studentID int64,
	page int,
	pageSize int,
	examName string,
	startDate string,
	endDate string,
) (*PageResult, error) {

	var (
		records []data_collection.AcademicRecord
		total   int64
	)

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		// ========= 1️⃣ 校验用户 =========
		var count int64
		if err := tx.Model(&users.User{}).
			Where("user_id = ?", studentID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("查询用户失败: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("学生不存在")
		}

		// ========= 2️⃣ 构建查询 =========
		query := tx.Model(&data_collection.AcademicRecord{}).
			Where("user_id = ?", studentID)

		// ========= 3️⃣ 模糊查询（考试名称）=========
		if examName != "" {
			query = query.Where("exam_name LIKE ?", "%"+examName+"%")
		}

		// ========= 4️⃣ 时间范围 =========
		if startDate != "" {
			t, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				return fmt.Errorf("start_date 格式错误，应为 yyyy-MM-dd")
			}
			query = query.Where("exam_date >= ?", t)
		}

		if endDate != "" {
			t, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				return fmt.Errorf("end_date 格式错误，应为 yyyy-MM-dd")
			}
			query = query.Where("exam_date <= ?", t)
		}

		// ========= 5️⃣ 统计总数 =========
		if err := query.Count(&total).Error; err != nil {
			return fmt.Errorf("统计总数失败: %w", err)
		}

		// ========= 6️⃣ 分页查询 =========
		offset := (page - 1) * pageSize

		if err := query.
			Order("exam_date DESC").
			Limit(pageSize).
			Offset(offset).
			Find(&records).Error; err != nil {
			return fmt.Errorf("查询数据失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &PageResult{
		List:     records,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
