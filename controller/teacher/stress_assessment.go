package teacher

import (
	"Monitoring-Pressure/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetStudentStressAssessmentListHandle 获取学生压力评估记录（分页 + 筛选 + 时间范围）
func GetStudentStressAssessmentListHandle(c *gin.Context) {

	// ========= 1️⃣ student_id =========
	studentIDStr := c.Query("student_id")
	if studentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id 不能为空"})
		return
	}

	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil || studentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id 格式错误"})
		return
	}

	// ========= 2️⃣ 分页 =========
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// ========= 3️⃣ 筛选 =========
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var warningStatus *int
	if ws := c.Query("warning_status"); ws != "" {
		v, err := strconv.Atoi(ws)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "warning_status 参数错误"})
			return
		}
		warningStatus = &v
	}

	// ========= 4️⃣ 调用 service =========
	result, err := service.GetStressAssessmentList(service.StressQuery{
		UserID:        &studentID,
		StartDate:     startDate,
		EndDate:       endDate,
		WarningStatus: warningStatus,
		Page:          page,
		PageSize:      pageSize,
	})

	if err != nil {
		log.Printf("获取学生[%d]压力评估记录失败: %v", studentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取压力评估记录失败",
		})
		return
	}

	// ========= 5️⃣ 返回 =========
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    result,
	})
}
