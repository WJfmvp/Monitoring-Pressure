package teacher

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetStudentInterventionRecordListHandle(c *gin.Context) {

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
	feedbackStatusStr := c.Query("feedback_status")

	// ========= 4️⃣ 调用 service =========
	result, err := service.GetStudentInterventionRecordsAdvanced(
		studentID,
		page,
		pageSize,
		startDate,
		endDate,
		feedbackStatusStr,
	)

	if err != nil {
		log.Printf("获取干预记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    result,
	})
}
