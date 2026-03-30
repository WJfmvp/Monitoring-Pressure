package models

import (
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/models/users"
	"fmt"
	"gorm.io/gorm"
)

// Migrate 执行数据库迁移，并且使用事务保证迁移的原子性
func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database instance is nil")
	}

	// 所有需要迁移的模型统一管理
	models := []interface{}{
		&users.User{},
		&data_collection.StudyBehaviorRecord{},
		&data_collection.AcademicImportRecord{},
		&data_collection.AcademicRecord{},
		&data_collection.QuestionnaireRecord{},
		&data_collection.PsychologicalSelfAssessment{},
		&data_collection.StressFeatureSnapshot{},
		&data_collection.StressAssessmentResult{},
		&data_collection.InterventionSuggestion{},
		&data_collection.StudentInterventionRecord{},
	}

	// 使用事务执行迁移
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, model := range models {
			if err := tx.AutoMigrate(model); err != nil {
				return fmt.Errorf("failed to migrate model %T: %w", model, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
