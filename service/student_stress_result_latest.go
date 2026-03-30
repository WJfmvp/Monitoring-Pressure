package service

import (
	"Monitoring-Pressure/dao/db"
	"errors"

	DataCollection "Monitoring-Pressure/models/data_collection"

	"gorm.io/gorm"
)

type StudentLatestStressResultResp struct {
	ID uint `json:"id"`

	UserID int64 `json:"user_id"`

	SnapshotID uint `json:"snapshot_id"`

	BehaviorScore      float64 `json:"behavior_score"`
	AcademicScore      float64 `json:"academic_score"`
	PsychologicalScore float64 `json:"psychological_score"`
	TotalScore         float64 `json:"total_score"`

	ModelPredictLevel  int     `json:"model_predict_level"`
	PredictProbability float64 `json:"predict_probability"`

	WarningStatus int    `json:"warning_status"`
	WarningTime   string `json:"warning_time"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetStudentLatestStressResult(userID int64) (*StudentLatestStressResultResp, error) {
	if userID <= 0 {
		return nil, errors.New("user_id 不合法")
	}

	var result DataCollection.StressAssessmentResult

	err := db.DB.
		Model(&DataCollection.StressAssessmentResult{}).
		Where("user_id = ?", userID).
		Order("id DESC").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var warningTimeStr string
	if result.WarningTime != nil {
		warningTimeStr = result.WarningTime.Format("2006-01-02 15:04:05")
	}

	resp := &StudentLatestStressResultResp{
		ID: result.ID,

		UserID: result.UserID,

		SnapshotID: result.SnapshotID,

		BehaviorScore:      result.BehaviorScore,
		AcademicScore:      result.AcademicScore,
		PsychologicalScore: result.PsychologicalScore,
		TotalScore:         result.TotalScore,

		ModelPredictLevel:  result.ModelPredictLevel,
		PredictProbability: result.PredictProbability,

		WarningStatus: result.WarningStatus,
		WarningTime:   warningTimeStr,

		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
