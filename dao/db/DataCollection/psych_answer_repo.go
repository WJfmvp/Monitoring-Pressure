package DataCollection

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/DataCollection"
	"gorm.io/gorm"
	"time"
)

type PsychAnswerRepo struct {
	db *gorm.DB
}

func NewPsychAnswerRepo() *PsychAnswerRepo {
	return &PsychAnswerRepo{db: db.DB}
}

// Create 新增问卷作答记录
func (repo *PsychAnswerRepo) Create(answer *DataCollection.PsychAnswer) error {
	answer.CreateTime = time.Now()
	answer.UpdateTime = time.Now()
	answer.SubmitTime = time.Now() // 提交时间默认当前
	return repo.db.Create(answer).Error
}

// GetByStudentAndQuestionnaire 根据学生ID和问卷ID查询作答记录
func (repo *PsychAnswerRepo) GetByStudentAndQuestionnaire(studentID string, questionnaireID int64) (*DataCollection.PsychAnswer, error) {
	var answer DataCollection.PsychAnswer
	err := repo.db.Where("student_id = ? AND questionnaire_id = ?", studentID, questionnaireID).First(&answer).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &answer, err
}

// UpdateSyncStatus 更新同步状态
func (repo *PsychAnswerRepo) UpdateSyncStatus(id int64, syncStatus int8, syncTime *time.Time) error {
	updateData := map[string]interface{}{
		"sync_status": syncStatus,
		"update_time": time.Now(),
	}
	if syncTime != nil {
		updateData["sync_time"] = syncTime
	}
	return repo.db.Model(&DataCollection.PsychAnswer{}).Where("id = ?", id).Updates(updateData).Error
}

// ListBySyncStatus 根据同步状态查询作答记录
func (repo *PsychAnswerRepo) ListBySyncStatus(syncStatus int8) ([]*DataCollection.PsychAnswer, error) {
	var answers []*DataCollection.PsychAnswer
	err := repo.db.Where("sync_status = ?", syncStatus).Find(&answers).Error
	return answers, err
}

// Delete 软删除作答记录
func (repo *PsychAnswerRepo) Delete(id int64) error {
	return repo.db.Delete(&DataCollection.PsychAnswer{}, id).Error
}
