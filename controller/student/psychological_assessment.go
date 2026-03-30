package student

import (
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// SubmitPsychologicalAssessmentHandle 提交心理问卷
func SubmitPsychologicalAssessmentHandle(c *gin.Context) {
	var assessment data_collection.PsychologicalSelfAssessment

	// 绑定前端传递的 JSON 数据到结构体
	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数格式错误，数据绑定失败",
		})
		return
	}

	// 数据验证：确保核心字段不为空
	if assessment.QuestionnaireID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "问卷ID不能为空",
		})
		return
	}

	if assessment.AssessDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "评估日期不能为空",
		})
		return
	}

	// 获取用户ID，通常是通过JWT解码获取
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}

	// 设置用户ID
	assessment.UserID = userID.(int64)
	assessment.CreatedAt = time.Now()
	assessment.UpdatedAt = time.Now()

	// 校验问卷分数：确保其合理性
	if assessment.QuestionnaireScore < 0 || assessment.QuestionnaireScore > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "问卷分数必须在0到100之间",
		})
		return
	}

	// 调用服务层处理并保存数据
	err := service.SubmitPsychologicalAssessment(assessment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("提交失败: %v", err),
		})
		return
	}

	// 提交成功
	c.JSON(http.StatusOK, gin.H{
		"message": "心理问卷提交成功",
		"data":    assessment,
	})
}

// GetPsychologicalAssessmentListHandle 获取心理自评问卷列表
func GetPsychologicalAssessmentListHandle(c *gin.Context) {
	// 获取用户ID，通常是通过 JWT 解码获取
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		log.Printf("无效请求：未提供用户ID，无法获取心理自评问卷列表")
		return
	}

	// 调用服务层获取该用户的所有心理评估记录
	assessments, err := service.GetPsychologicalAssessmentsByUserID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取心理评估记录失败",
		})
		log.Printf("获取用户 %d 心理评估记录失败: %v", userID, err)
		return
	}

	// 如果没有找到记录，返回空列表
	if len(assessments) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有找到任何心理评估记录",
			"data":    assessments,
		})
		return
	}

	// 返回获取到的数据
	c.JSON(http.StatusOK, gin.H{
		"message": "获取心理评估记录成功",
		"data":    assessments,
	})
}
