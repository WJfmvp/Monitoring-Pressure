package DataCollection

import "time"

// StressFeatureSnapshot 把三类数据汇总后的“训练样本表”
// 把行为 + 成绩 + 心理聚合成一条样本
type StressFeatureSnapshot struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	UserID int64 `json:"user_id" gorm:"not null;index"`

	SnapshotDate time.Time `json:"snapshot_date" gorm:"not null;type:date"`

	// ===== 学习行为特征 =====
	StudyDuration      float64 `json:"study_duration" gorm:"type:decimal(5,2);default:0"`
	HomeworkSubmitRate float64 `json:"homework_submit_rate" gorm:"type:decimal(5,2);default:0"`
	StudyFrequency     int     `json:"study_frequency" gorm:"default:0"`
	AbsenceCount       int     `json:"absence_count" gorm:"default:0"`

	// ===== 学业成绩特征 =====
	ExamScore        float64 `json:"exam_score" gorm:"type:decimal(6,2);default:0"`
	ScoreFluctuation float64 `json:"score_fluctuation" gorm:"type:decimal(6,2);default:0"`
	ClassRank        int     `json:"class_rank" gorm:"default:0"`

	// ===== 心理特征 =====
	AnxietyLevel       int `json:"anxiety_level" gorm:"default:0"`
	LearningMotivation int `json:"learning_motivation" gorm:"default:0"`
	EmotionalState     int `json:"emotional_state" gorm:"default:0"`
	SleepQuality       int `json:"sleep_quality" gorm:"default:0"`

	// ===== 综合维度得分（用于规则评估）=====
	BehaviorScore      float64 `json:"behavior_score" gorm:"type:decimal(6,2);default:0"`
	AcademicScore      float64 `json:"academic_score" gorm:"type:decimal(6,2);default:0"`
	PsychologicalScore float64 `json:"psychological_score" gorm:"type:decimal(6,2);default:0"`

	// 标签：0低 1中 2高
	StressLevel int `json:"stress_level" gorm:"not null;default:0"`

	FeatureVersion string `json:"feature_version" gorm:"type:varchar(20);default:'v1'"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (StressFeatureSnapshot) TableName() string {
	return "stress_feature_snapshots"
}
