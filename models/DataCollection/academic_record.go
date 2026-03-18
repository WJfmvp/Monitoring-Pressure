package DataCollection

import "time"

// AcademicRecord 存每个学生每次考试成绩
// 来源：Excel 解析后批量入库
type AcademicRecord struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	UserID int64 `json:"user_id" gorm:"not null;index"`

	// 对应哪次Excel导入
	ImportID uint `json:"import_id" gorm:"index"`

	ExamName string    `json:"exam_name" gorm:"type:varchar(100);not null"` // 考试名称
	Term     string    `json:"term" gorm:"type:varchar(50)"`                // 学期
	ExamDate time.Time `json:"exam_date" gorm:"type:date"`                  // 考试日期

	// 原始成绩（导入值）
	RawScore float64 `json:"raw_score" gorm:"type:decimal(6,2);default:0"`

	// 清洗后的成绩
	ExamScore        float64 `json:"exam_score" gorm:"type:decimal(6,2);default:0"`
	AverageScore     float64 `json:"average_score" gorm:"type:decimal(6,2);default:0"`
	ScoreFluctuation float64 `json:"score_fluctuation" gorm:"type:decimal(6,2);default:0"`

	// 排名
	ClassRank int `json:"class_rank" gorm:"default:0"`
	GradeRank int `json:"grade_rank" gorm:"default:0"`

	// 扩展指标
	FailCount       int     `json:"fail_count" gorm:"default:0"`                          // 挂科数
	ScoreChangeRate float64 `json:"score_change_rate" gorm:"type:decimal(6,2);default:0"` // 成绩变化率

	// 数据来源：1手动录入 2Excel导入
	SourceType int `json:"source_type" gorm:"not null;default:2"`

	Remark string `json:"remark" gorm:"type:varchar(255)"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AcademicRecord) TableName() string {
	return "academic_records"
}
