package data_collection

import (
	"Monitoring-Pressure/models/users"
	"time"
)

// PsychologicalSelfAssessment 存根据问卷整理出来的心理特征结果
// 来源：前端问卷提交后写入
type PsychologicalSelfAssessment struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	UserID int64      `json:"user_id" gorm:"not null;index"`
	User   users.User `gorm:"foreignKey:UserID;references:UserID"`

	QuestionnaireID uint                `json:"questionnaire_id" gorm:"index"` // 对应问卷记录
	Questionnaire   QuestionnaireRecord `gorm:"foreignKey:QuestionnaireID;references:ID"`

	AssessDate time.Time `json:"assess_date" gorm:"not null;type:date"`

	// 核心指标
	AnxietyLevel       int `json:"anxiety_level" gorm:"default:0"`       // 焦虑程度
	LearningMotivation int `json:"learning_motivation" gorm:"default:0"` // 学习动力
	EmotionalState     int `json:"emotional_state" gorm:"default:0"`     // 情绪状态

	// 扩展指标
	StressPerception int `json:"stress_perception" gorm:"default:0"` // 自我压力感知
	SleepQuality     int `json:"sleep_quality" gorm:"default:0"`     // 睡眠质量
	FatigueLevel     int `json:"fatigue_level" gorm:"default:0"`     // 疲劳程度

	QuestionnaireScore float64 `json:"questionnaire_score" gorm:"type:decimal(6,2);default:0"`

	Remark string `json:"remark" gorm:"type:varchar(255)"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PsychologicalSelfAssessment) TableName() string {
	return "psychological_self_assessments"
}
