package data_collection

import (
	"Monitoring-Pressure/models/users"
	"time"
)

// StressAssessmentResult 存最终压力结果、模型预测、预警结果
// 来源：Python 模型预测后写回数据库
type StressAssessmentResult struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	UserID int64      `json:"user_id" gorm:"not null;index"`
	User   users.User `gorm:"foreignKey:UserID;references:UserID"`

	SnapshotID uint `json:"snapshot_id" gorm:"not null;index"`

	// 规则评分
	BehaviorScore      float64 `json:"behavior_score" gorm:"type:decimal(6,2);default:0"`
	AcademicScore      float64 `json:"academic_score" gorm:"type:decimal(6,2);default:0"`
	PsychologicalScore float64 `json:"psychological_score" gorm:"type:decimal(6,2);default:0"`
	TotalScore         float64 `json:"total_score" gorm:"type:decimal(6,2);default:0"`

	// 模型结果
	ModelPredictLevel  int     `json:"model_predict_level" gorm:"default:0"` // 0低 1中 2高
	PredictProbability float64 `json:"predict_probability" gorm:"type:decimal(6,4);default:0"`

	// 预警状态：0正常 1中预警 2高预警
	WarningStatus int `json:"warning_status" gorm:"default:0"`

	WarningTime *time.Time `json:"warning_time"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (StressAssessmentResult) TableName() string {
	return "stress_assessment_results"
}
