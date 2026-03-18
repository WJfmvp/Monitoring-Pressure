package DataCollection

import (
	"Monitoring-Pressure/models/users"
	"time"
)

// StudyBehaviorRecord 学习行为
// 数据来源：学生填写 教师录入 系统导入
type StudyBehaviorRecord struct {
	ID     uint  `gorm:"primaryKey"`
	UserID int64 `gorm:"not null"`

	RecordDate time.Time `gorm:"not null"`

	// 作业提交率
	StudyFrequency int `gorm:"default:0"` // 学习频率

	// ===== 学业成绩特征 =====
	ExamScore        float64 `gorm:"default:0"`
	ScoreFluctuation float64 `gorm:"default:0"`
	ClassRank        int     `gorm:"default:0"`

	// ===== 心理特征 =====
	AnxietyLevel       int `gorm:"default:0"`
	LearningMotivation int `gorm:"default:0"`
	EmotionalState     int `gorm:"default:0"`

	// ===== 标签（关键！！）=====
	StressLevel int `gorm:"default:0"` // 0低 1中 2高（训练用标签）

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User users.User `gorm:"foreignKey:StudentID"`
}
