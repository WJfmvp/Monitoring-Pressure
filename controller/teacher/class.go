package teacher

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Monitoring-Pressure/service"
)

// GetClassStressOverviewHandle 班级压力总览（聚合每个学生最新一次评估）
// GET /teacher/class/overview
func GetClassStressOverviewHandle(c *gin.Context) {
	overview, err := service.GetClassStressOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取班级总览失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": overview,
	})
}

// GetStudentDiagnosisHandle 单个学生的压力诊断书
// GET /teacher/student/diagnosis?student_id=X
func GetStudentDiagnosisHandle(c *gin.Context) {
	studentIDStr := c.Query("student_id")
	if studentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "student_id 不能为空",
		})
		return
	}
	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil || studentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "student_id 格式错误",
		})
		return
	}

	diag, err := service.GetStudentDiagnosis(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": diag,
	})
}
