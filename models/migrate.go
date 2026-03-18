package models

import (
	"Monitoring-Pressure/models/DataCollection"
	"Monitoring-Pressure/models/users"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&users.User{},
		&DataCollection.StudyBehaviorRecord{},
		&DataCollection.AcademicImportRecord{},
		&DataCollection.AcademicRecord{},
		&DataCollection.QuestionnaireRecord{},
		&DataCollection.PsychologicalSelfAssessment{},
		&DataCollection.StressFeatureSnapshot{},
		&DataCollection.StressAssessmentResult{},
		&DataCollection.InterventionSuggestion{},
		&DataCollection.StudentInterventionRecord{},
	)
}
