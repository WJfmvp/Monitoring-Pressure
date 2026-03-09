package DataCollection

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/DataCollection"
	"gorm.io/gorm"
	"time"
)

type PsychQuestionnaireTemplateRepo struct {
	db *gorm.DB
}

func NewPsychQuestionnaireTemplateRepo() *PsychQuestionnaireTemplateRepo {
	return &PsychQuestionnaireTemplateRepo{db: db.DB}
}

// Create 新增问卷模板
func (repo *PsychQuestionnaireTemplateRepo) Create(template *DataCollection.PsychQuestionnaireTemplate) error {
	template.CreateTime = time.Now()
	template.UpdateTime = time.Now()
	return repo.db.Create(template).Error
}

// GetByID 根据ID查询模板
func (repo *PsychQuestionnaireTemplateRepo) GetByID(id int64) (*DataCollection.PsychQuestionnaireTemplate, error) {
	var template DataCollection.PsychQuestionnaireTemplate
	err := repo.db.Where("id = ?", id).First(&template).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &template, err
}

// ListAll 查询所有有效模板（未过期）
func (repo *PsychQuestionnaireTemplateRepo) ListAll() ([]*DataCollection.PsychQuestionnaireTemplate, error) {
	var templates []*DataCollection.PsychQuestionnaireTemplate
	now := time.Now()
	err := repo.db.Where("valid_end_time >= ?", now).Find(&templates).Error
	return templates, err
}

// Update 更新模板
func (repo *PsychQuestionnaireTemplateRepo) Update(template *DataCollection.PsychQuestionnaireTemplate) error {
	template.UpdateTime = time.Now()
	return repo.db.Model(template).Select(
		"questionnaire_name", "dimension", "is_required",
		"valid_start_time", "valid_end_time", "update_time",
	).Where("id = ?", template.ID).Updates(template).Error
}

// Delete 软删除模板
func (repo *PsychQuestionnaireTemplateRepo) Delete(id int64) error {
	return repo.db.Delete(&DataCollection.PsychQuestionnaireTemplate{}, id).Error
}
