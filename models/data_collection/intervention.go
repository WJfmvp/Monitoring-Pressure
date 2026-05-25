package data_collection

import (
	"Monitoring-Pressure/models/users"
	"time"
)

// InterventionSuggestion 系统里的建议模板库
// 来源:管理员预先配置
type InterventionSuggestion struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	Level int `json:"level" gorm:"not null;default:0"` // 0低 1中 2高

	Title   string `json:"title" gorm:"type:varchar(100);not null"`
	Content string `json:"content" gorm:"type:text;not null"`

	Category string `json:"category" gorm:"type:varchar(50)"` // 学习/心理/生活

	MinScore float64 `json:"min_score" gorm:"type:decimal(6,2);default:0"`
	MaxScore float64 `json:"max_score" gorm:"type:decimal(6,2);default:0"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (InterventionSuggestion) TableName() string {
	return "intervention_suggestions"
}

// StudentInterventionRecord 记录某个学生收到过什么建议、反馈如何
// 来源:推进建议生成记录
type StudentInterventionRecord struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	UserID int64      `json:"user_id" gorm:"not null;index"`
	User   users.User `gorm:"foreignKey:UserID;references:UserID"`

	ResultID uint                   `json:"result_id" gorm:"not null;index"`
	Result   StressAssessmentResult `json:"result" gorm:"foreignKey:ResultID;references:ID"`

	SuggestionID uint                   `json:"suggestion_id" gorm:"not null;index"`
	Suggestion   InterventionSuggestion `json:"suggestion" gorm:"foreignKey:SuggestionID;references:ID"`

	PushTime time.Time `json:"push_time" gorm:"not null"`

	Feedback       string `json:"feedback" gorm:"type:varchar(255)"`
	FeedbackStatus int    `json:"feedback_status" gorm:"default:0"` // 0未反馈 1有效 2无效

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (StudentInterventionRecord) TableName() string {
	return "student_intervention_records"
}
