package service

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/models/users"
	"fmt"
	"gorm.io/gorm"
)

// SubmitPsychologicalAssessment 提交心理自评问卷
func SubmitPsychologicalAssessment(assessment data_collection.PsychologicalSelfAssessment) error {
	// 使用事务来确保数据一致性
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var userCount int64
		err := tx.Model(&users.User{}).Where("user_id = ?", assessment.UserID).Count(&userCount).Error
		if err != nil {
			return fmt.Errorf("查询用户信息失败: %w", err)
		}
		if userCount == 0 {
			return fmt.Errorf("用户不存在")
		}

		// 保存数据到数据库
		if err := tx.Create(&assessment).Error; err != nil {
			return fmt.Errorf("保存心理问卷失败: %w", err)
		}

		// 如果有其他需要更新的逻辑（例如记录、提醒等），可以在此扩展

		return nil
	})
}

// GetPsychologicalAssessmentsByUserID 获取用户的所有心理自评记录
func GetPsychologicalAssessmentsByUserID(userID int64) ([]data_collection.PsychologicalSelfAssessment, error) {
	var assessments []data_collection.PsychologicalSelfAssessment

	// 查询指定用户的所有心理自评问卷记录
	err := db.DB.Where("user_id = ?", userID).Find(&assessments).Error
	if err != nil {
		// 记录数据库查询错误
		return nil, fmt.Errorf("查询心理评估记录失败: %w", err)
	}

	// 如果没有记录，返回空切片
	if len(assessments) == 0 {
		return assessments, nil
	}

	// 返回查询结果
	return assessments, nil
}
