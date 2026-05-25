package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/models/users"

	"gorm.io/gorm"
)

// StressPredictRequest 触发一次压力预测所需的特征字段
// 字段语义与 Python 模型输入完全一致，由调用方（管理员/教师）填写
type StressPredictRequest struct {
	UserID int64 `json:"user_id" binding:"required"`

	SnapshotDate string `json:"snapshot_date"` // 可选 yyyy-MM-dd，缺省取今天

	// ===== 学习行为特征 =====
	StudyDuration      float64 `json:"study_duration"`
	HomeworkSubmitRate float64 `json:"homework_submit_rate"`
	StudyFrequency     int     `json:"study_frequency"`
	AbsenceCount       int     `json:"absence_count"`

	// ===== 学业成绩特征 =====
	ExamScore        float64 `json:"exam_score"`
	ScoreFluctuation float64 `json:"score_fluctuation"`
	ClassRank        int     `json:"class_rank"`

	// ===== 心理特征 =====
	AnxietyLevel       int `json:"anxiety_level"`
	LearningMotivation int `json:"learning_motivation"`
	EmotionalState     int `json:"emotional_state"`
	SleepQuality       int `json:"sleep_quality"`

	// ===== 综合维度得分 =====
	BehaviorScore      float64 `json:"behavior_score"`
	AcademicScore      float64 `json:"academic_score"`
	PsychologicalScore float64 `json:"psychological_score"`
}

// StressPredictResp 给前端返回的最终结果
type StressPredictResp struct {
	ResultID   uint `json:"result_id"`
	SnapshotID uint `json:"snapshot_id"`

	UserID int64 `json:"user_id"`

	BehaviorScore      float64 `json:"behavior_score"`
	AcademicScore      float64 `json:"academic_score"`
	PsychologicalScore float64 `json:"psychological_score"`
	TotalScore         float64 `json:"total_score"` // 与 Python 的 risk_index 对齐

	ModelPredictLevel  int     `json:"model_predict_level"`
	PredictProbability float64 `json:"predict_probability"`

	WarningStatus int    `json:"warning_status"` // 与 Python 的 final_level 对齐
	WarningTime   string `json:"warning_time"`

	MentalRisk      float64 `json:"mental_risk"`
	BehaviorRisk    float64 `json:"behavior_risk"`
	HighAnxietyFlag int     `json:"high_anxiety_flag"`

	CreatedAt string `json:"created_at"`
}

// PredictAndSaveStress 完整压力检测流程：建快照 → 调 Python → 写结果
func PredictAndSaveStress(req StressPredictRequest) (*StressPredictResp, error) {
	if req.UserID <= 0 {
		return nil, errors.New("user_id 不合法")
	}

	if err := validatePredictRequest(req); err != nil {
		return nil, err
	}

	// 解析快照日期（缺省今天）
	snapshotDate := time.Now()
	if req.SnapshotDate != "" {
		t, err := time.Parse("2006-01-02", req.SnapshotDate)
		if err != nil {
			return nil, fmt.Errorf("snapshot_date 格式错误，应为 yyyy-MM-dd")
		}
		snapshotDate = t
	}

	// 1. 校验用户存在 + 写入快照（事务内，避免 HTTP 长占连接）
	snapshot := data_collection.StressFeatureSnapshot{
		UserID:             req.UserID,
		SnapshotDate:       snapshotDate,
		StudyDuration:      req.StudyDuration,
		HomeworkSubmitRate: req.HomeworkSubmitRate,
		StudyFrequency:     req.StudyFrequency,
		AbsenceCount:       req.AbsenceCount,
		ExamScore:          req.ExamScore,
		ScoreFluctuation:   req.ScoreFluctuation,
		ClassRank:          req.ClassRank,
		AnxietyLevel:       req.AnxietyLevel,
		LearningMotivation: req.LearningMotivation,
		EmotionalState:     req.EmotionalState,
		SleepQuality:       req.SleepQuality,
		BehaviorScore:      req.BehaviorScore,
		AcademicScore:      req.AcademicScore,
		PsychologicalScore: req.PsychologicalScore,
		FeatureVersion:     "v1",
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var userCount int64
		if err := tx.Model(&users.User{}).
			Where("user_id = ?", req.UserID).
			Count(&userCount).Error; err != nil {
			return fmt.Errorf("校验用户失败: %w", err)
		}
		if userCount == 0 {
			return fmt.Errorf("用户不存在")
		}

		if err := tx.Create(&snapshot).Error; err != nil {
			return fmt.Errorf("保存特征快照失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 2. 调 Python 随机森林服务（事务外，避免长锁）
	mlReq := MLPredictRequest{
		StudentID:          strconv.FormatInt(req.UserID, 10),
		StudyDuration:      req.StudyDuration,
		HomeworkSubmitRate: req.HomeworkSubmitRate,
		StudyFrequency:     req.StudyFrequency,
		AbsenceCount:       req.AbsenceCount,
		ExamScore:          req.ExamScore,
		ScoreFluctuation:   req.ScoreFluctuation,
		ClassRank:          req.ClassRank,
		AnxietyLevel:       req.AnxietyLevel,
		LearningMotivation: req.LearningMotivation,
		EmotionalState:     req.EmotionalState,
		SleepQuality:       req.SleepQuality,
		BehaviorScore:      req.BehaviorScore,
		AcademicScore:      req.AcademicScore,
		PsychologicalScore: req.PsychologicalScore,
	}

	mlData, err := CallMLPredict(mlReq)
	if err != nil {
		return nil, err
	}

	// 3. 写入压力评估结果
	var warningTime *time.Time
	if mlData.FinalLevel > 0 {
		now := time.Now()
		warningTime = &now
	}

	result := data_collection.StressAssessmentResult{
		UserID:             req.UserID,
		SnapshotID:         snapshot.ID,
		BehaviorScore:      req.BehaviorScore,
		AcademicScore:      req.AcademicScore,
		PsychologicalScore: req.PsychologicalScore,
		TotalScore:         mlData.RiskIndex, // 与 Python risk_index 对齐
		ModelPredictLevel:  mlData.ModelPredictLevel,
		PredictProbability: mlData.PredictProbability,
		WarningStatus:      mlData.FinalLevel, // 规则后调整的最终等级
		WarningTime:        warningTime,
	}

	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&result).Error; err != nil {
			return fmt.Errorf("保存压力评估结果失败: %w", err)
		}
		if err := GenerateInterventionRecordsForResult(tx, result); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// 4. 整理响应
	resp := &StressPredictResp{
		ResultID:           result.ID,
		SnapshotID:         snapshot.ID,
		UserID:             result.UserID,
		BehaviorScore:      result.BehaviorScore,
		AcademicScore:      result.AcademicScore,
		PsychologicalScore: result.PsychologicalScore,
		TotalScore:         result.TotalScore,
		ModelPredictLevel:  result.ModelPredictLevel,
		PredictProbability: result.PredictProbability,
		WarningStatus:      result.WarningStatus,
		MentalRisk:         mlData.MentalRisk,
		BehaviorRisk:       mlData.BehaviorRisk,
		HighAnxietyFlag:    mlData.HighAnxietyFlag,
		CreatedAt:          result.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if warningTime != nil {
		resp.WarningTime = warningTime.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}

// validatePredictRequest 与 Python pydantic schema 中的范围保持一致
func validatePredictRequest(r StressPredictRequest) error {
	if r.StudyDuration < 0 || r.StudyDuration > 24 {
		return errors.New("study_duration 必须在 0~24 之间")
	}
	if r.HomeworkSubmitRate < 0 || r.HomeworkSubmitRate > 1 {
		return errors.New("homework_submit_rate 必须在 0~1 之间")
	}
	if r.StudyFrequency < 0 || r.StudyFrequency > 7 {
		return errors.New("study_frequency 必须在 0~7 之间")
	}
	if r.AbsenceCount < 0 || r.AbsenceCount > 365 {
		return errors.New("absence_count 必须在 0~365 之间")
	}
	if r.ExamScore < 0 || r.ExamScore > 100 {
		return errors.New("exam_score 必须在 0~100 之间")
	}
	if r.ScoreFluctuation < 0 || r.ScoreFluctuation > 100 {
		return errors.New("score_fluctuation 必须在 0~100 之间")
	}
	if r.ClassRank < 1 || r.ClassRank > 100000 {
		return errors.New("class_rank 必须在 1~100000 之间")
	}
	if r.AnxietyLevel < 0 || r.AnxietyLevel > 10 {
		return errors.New("anxiety_level 必须在 0~10 之间")
	}
	if r.LearningMotivation < 0 || r.LearningMotivation > 10 {
		return errors.New("learning_motivation 必须在 0~10 之间")
	}
	if r.EmotionalState < 0 || r.EmotionalState > 10 {
		return errors.New("emotional_state 必须在 0~10 之间")
	}
	if r.SleepQuality < 0 || r.SleepQuality > 10 {
		return errors.New("sleep_quality 必须在 0~10 之间")
	}
	if r.BehaviorScore < 0 || r.BehaviorScore > 100 {
		return errors.New("behavior_score 必须在 0~100 之间")
	}
	if r.AcademicScore < 0 || r.AcademicScore > 100 {
		return errors.New("academic_score 必须在 0~100 之间")
	}
	if r.PsychologicalScore < 0 || r.PsychologicalScore > 100 {
		return errors.New("psychological_score 必须在 0~100 之间")
	}
	return nil
}
