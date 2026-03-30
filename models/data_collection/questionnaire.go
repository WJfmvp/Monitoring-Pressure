package data_collection

import (
	"Monitoring-Pressure/models/users"
	"time"
)

// QuestionnaireRecord 记录学生填写了一次问卷
// 来源：前端提交问卷时先保存一条记录
type QuestionnaireRecord struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	UserID int64      `json:"user_id" gorm:"not null;index"`
	User   users.User `gorm:"foreignKey:UserID;references:UserID"`

	Title string `json:"title" gorm:"type:varchar(100);not null"` // 问卷名称
	Type  string `json:"type" gorm:"type:varchar(50)"`            // 问卷类型

	TotalScore float64 `json:"total_score" gorm:"type:decimal(6,2);default:0"` // 总分

	Status int `json:"status" gorm:"not null;default:1"` // 1有效 0无效

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (QuestionnaireRecord) TableName() string {
	return "questionnaire_records"
}
