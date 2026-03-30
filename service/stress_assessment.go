package service

import (
	"fmt"
	"time"

	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/data_collection"
	"gorm.io/gorm"
)

type StressQuery struct {
	UserID        *int64
	StartDate     string
	EndDate       string
	WarningStatus *int
	Page          int
	PageSize      int
}

// GetStressAssessmentList 获取压力评估记录（高级查询）
func GetStressAssessmentList(queryParam StressQuery) (*PageResult, error) {

	var (
		list  []data_collection.StressAssessmentResult
		total int64
	)

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		query := tx.Model(&data_collection.StressAssessmentResult{})

		// ========= 条件拼接 =========

		if queryParam.UserID != nil {
			query = query.Where("user_id = ?", *queryParam.UserID)
		}

		if queryParam.StartDate != "" {
			t, err := time.Parse("2006-01-02", queryParam.StartDate)
			if err != nil {
				return fmt.Errorf("start_date 格式错误")
			}
			query = query.Where("created_at >= ?", t)
		}

		if queryParam.EndDate != "" {
			t, err := time.Parse("2006-01-02", queryParam.EndDate)
			if err != nil {
				return fmt.Errorf("end_date 格式错误")
			}
			query = query.Where("created_at <= ?", t)
		}

		if queryParam.WarningStatus != nil {
			query = query.Where("warning_status = ?", *queryParam.WarningStatus)
		}

		// ========= count =========
		if err := query.Count(&total).Error; err != nil {
			return err
		}

		// ========= 分页 =========
		offset := (queryParam.Page - 1) * queryParam.PageSize

		if err := query.
			Order("created_at DESC").
			Limit(queryParam.PageSize).
			Offset(offset).
			Find(&list).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &PageResult{
		List:     list,
		Total:    total,
		Page:     queryParam.Page,
		PageSize: queryParam.PageSize,
	}, nil
}
