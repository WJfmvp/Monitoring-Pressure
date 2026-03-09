package DataCollection

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/DataCollection"
	"gorm.io/gorm"
	"time"
)

type AcademicScoreRepo struct {
	db *gorm.DB
}

func NewAcademicScoreRepo() *AcademicScoreRepo {
	return &AcademicScoreRepo{db: db.DB}
}

// Create 新增单条成绩数据
func (repo *AcademicScoreRepo) Create(score *DataCollection.AcademicScore) error {
	score.CreateTime = time.Now()
	score.UpdateTime = time.Now()
	return repo.db.Create(score).Error
}

// GetByID 根据ID查询
func (repo *AcademicScoreRepo) GetByID(id int64) (*DataCollection.AcademicScore, error) {
	var score DataCollection.AcademicScore
	err := repo.db.Where("id = ?", id).First(&score).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &score, err
}

// ListByStudentID 根据学生ID查询成绩列表
func (repo *AcademicScoreRepo) ListByStudentID(studentID string) ([]*DataCollection.AcademicScore, error) {
	var scores []*DataCollection.AcademicScore
	err := repo.db.Where("student_id = ?", studentID).Find(&scores).Error
	return scores, err
}

// Update 更新成绩数据
func (repo *AcademicScoreRepo) Update(score *DataCollection.AcademicScore) error {
	score.UpdateTime = time.Now()
	return repo.db.Model(score).Select(
		"exam_name", "subject_name", "score", "grade",
		"ranking", "data_source", "update_time",
	).Where("id = ?", score.ID).Updates(score).Error
}

// Delete 软删除
func (repo *AcademicScoreRepo) Delete(id int64) error {
	return repo.db.Delete(&DataCollection.AcademicScore{}, id).Error
}
