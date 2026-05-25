package student

import (
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/service"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// submitPsychologicalReq 用 string 接收日期，避免 time.Time 默认要求 RFC3339
type submitPsychologicalReq struct {
	QuestionnaireID uint   `json:"questionnaire_id" binding:"required"`
	AssessDate      string `json:"assess_date" binding:"required"` // yyyy-MM-dd

	AnxietyLevel       int `json:"anxiety_level"`
	LearningMotivation int `json:"learning_motivation"`
	EmotionalState     int `json:"emotional_state"`

	StressPerception int `json:"stress_perception"`
	SleepQuality     int `json:"sleep_quality"`
	FatigueLevel     int `json:"fatigue_level"`

	QuestionnaireScore float64 `json:"questionnaire_score"`
	Remark             string  `json:"remark"`
}

// computePsychologicalScore 服务端兜底总分计算（与前端公式一致）
// 6 项指标各 0-10，正向指标越高压力越大，反向指标越高压力越小。归一到 0-100。
func computePsychologicalScore(r submitPsychologicalReq) float64 {
	clamp := func(v int) int {
		if v < 0 {
			return 0
		}
		if v > 10 {
			return 10
		}
		return v
	}
	positive := clamp(r.AnxietyLevel) + clamp(r.StressPerception) + clamp(r.FatigueLevel)
	inverted := (10 - clamp(r.LearningMotivation)) + (10 - clamp(r.EmotionalState)) + (10 - clamp(r.SleepQuality))
	raw := float64(positive + inverted) // 0..60
	score := raw / 60.0 * 100.0
	// 两位小数
	return float64(int(score*100+0.5)) / 100.0
}

// SubmitPsychologicalAssessmentHandle 提交心理问卷
func SubmitPsychologicalAssessmentHandle(c *gin.Context) {
	var req submitPsychologicalReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("请求参数格式错误: %v", err),
		})
		return
	}

	assessDate, err := time.Parse("2006-01-02", req.AssessDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "评估日期格式错误，应为 yyyy-MM-dd",
		})
		return
	}

	// 若前端没传或传了非法值，服务端兜底重算
	score := req.QuestionnaireScore
	if score <= 0 || score > 100 {
		score = computePsychologicalScore(req)
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}

	assessment := data_collection.PsychologicalSelfAssessment{
		UserID:             userID,
		QuestionnaireID:    req.QuestionnaireID,
		AssessDate:         assessDate,
		AnxietyLevel:       req.AnxietyLevel,
		LearningMotivation: req.LearningMotivation,
		EmotionalState:     req.EmotionalState,
		StressPerception:   req.StressPerception,
		SleepQuality:       req.SleepQuality,
		FatigueLevel:       req.FatigueLevel,
		QuestionnaireScore: score,
		Remark:             req.Remark,
	}

	if err := service.SubmitPsychologicalAssessment(assessment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("提交失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "心理问卷提交成功",
		"data":    assessment,
	})
}

// GetPsychologicalAssessmentListHandle 获取心理自评问卷列表
func GetPsychologicalAssessmentListHandle(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}

	assessments, err := service.GetPsychologicalAssessmentsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取心理评估记录失败",
		})
		log.Printf("获取用户 %d 心理评估记录失败: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取心理评估记录成功",
		"data":    assessments,
	})
}
