package DataCollection

import (
	"gorm.io/gorm"
	"time"
)

// LearningBehavior 学习行为表
type LearningBehavior struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	StudentID        string    `gorm:"column:student_id;type:varchar(32);not null;comment:学生唯一标识" json:"student_id"`
	CourseName       string    `gorm:"column:course_name;type:varchar(64);not null;comment:课程名称" json:"course_name"`
	LearningDate     time.Time `gorm:"column:learning_date;type:date;not null;comment:学习日期" json:"learning_date"`
	LearningDuration int       `gorm:"column:learning_duration;type:int;not null;default:0;check:learning_duration >= 0;comment:学习时长（分钟）" json:"learning_duration"`
	HomeworkStatus   int8      `gorm:"column:homework_status;type:tinyint;not null;comment:作业提交状态（0未提交/1已提交）" json:"homework_status"`
	HomeworkQuality  *int8     `gorm:"column:homework_quality;type:tinyint;comment:作业质量（1优/2良/3中/4差）" json:"homework_quality,omitempty"`
	DataSource       int8      `gorm:"column:data_source;type:tinyint;not null;comment:数据来源（1手动/2教务系统）" json:"data_source"`
	CreateBy         string    `gorm:"column:create_by;type:varchar(32);not null;comment:录入人ID" json:"create_by"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:录入时间" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	// 软删除（可选，根据业务需求添加）
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (lb *LearningBehavior) TableName() string {
	return "t_learning_behavior"
}

// AcademicScore 学业成绩表
type AcademicScore struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	StudentID   string         `gorm:"column:student_id;type:varchar(32);not null;comment:学生唯一标识" json:"student_id"`
	ExamName    string         `gorm:"column:exam_name;type:varchar(64);not null;comment:考试名称（如期中/期末）" json:"exam_name"`
	SubjectName string         `gorm:"column:subject_name;type:varchar(32);not null;comment:科目名称" json:"subject_name"`
	Score       float64        `gorm:"column:score;type:decimal(5,1);not null;check:score >= 0 and score <= 100;comment:分数" json:"score"`
	Grade       *string        `gorm:"column:grade;type:varchar(8);comment:等级（A/B/C/D）" json:"grade,omitempty"`
	ExamTime    time.Time      `gorm:"column:exam_time;type:date;not null;comment:考试时间" json:"exam_time"`
	Ranking     *int           `gorm:"column:ranking;type:int;comment:班级排名" json:"ranking,omitempty"`
	DataSource  int8           `gorm:"column:data_source;type:tinyint;not null;comment:数据来源（1手动/2教务系统）" json:"data_source"`
	CreateBy    string         `gorm:"column:create_by;type:varchar(32);not null;comment:录入人ID" json:"create_by"`
	CreateTime  time.Time      `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:录入时间" json:"create_time"`
	UpdateTime  time.Time      `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (as *AcademicScore) TableName() string {
	return "t_academic_score"
}

// PsychQuestionnaireTemplate 心理问卷模板表
type PsychQuestionnaireTemplate struct {
	ID                int64          `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	QuestionnaireName string         `gorm:"column:questionnaire_name;type:varchar(64);not null;comment:问卷名称" json:"questionnaire_name"`
	Dimension         string         `gorm:"column:dimension;type:varchar(256);not null;comment:问卷维度（逗号分隔，如焦虑、抑郁）" json:"dimension"`
	IsRequired        int8           `gorm:"column:is_required;type:tinyint;not null;comment:是否必填（0否/1是）" json:"is_required"`
	ValidStartTime    time.Time      `gorm:"column:valid_start_time;type:datetime;not null;comment:有效开始时间" json:"valid_start_time"`
	ValidEndTime      time.Time      `gorm:"column:valid_end_time;type:datetime;not null;comment:有效结束时间" json:"valid_end_time"`
	CreateTime        time.Time      `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime        time.Time      `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (pqt *PsychQuestionnaireTemplate) TableName() string {
	return "t_psych_questionnaire_template"
}

// PsychAnswer 问卷作答表
type PsychAnswer struct {
	ID              int64          `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	StudentID       string         `gorm:"column:student_id;type:varchar(32);not null;comment:学生唯一标识" json:"student_id"`
	QuestionnaireID int64          `gorm:"column:questionnaire_id;type:bigint;not null;comment:问卷模板ID" json:"questionnaire_id"`
	AnswerContent   string         `gorm:"column:answer_content;type:json;not null;comment:作答内容（题目ID-答案映射）" json:"answer_content"`
	TotalScore      float64        `gorm:"column:total_score;type:decimal(8,2);not null;comment:问卷总分" json:"total_score"`
	SyncStatus      int8           `gorm:"column:sync_status;type:tinyint;not null;comment:同步状态（0待同步/1同步中/2成功/3失败）" json:"sync_status"`
	SubmitTime      time.Time      `gorm:"column:submit_time;type:datetime;not null;comment:提交时间" json:"submit_time"`
	SyncTime        *time.Time     `gorm:"column:sync_time;type:datetime;comment:同步完成时间" json:"sync_time,omitempty"`
	CreateTime      time.Time      `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime      time.Time      `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (pa *PsychAnswer) TableName() string {
	return "t_psych_answer"
}
