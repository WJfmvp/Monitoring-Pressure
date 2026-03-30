package admin

import (
	"Monitoring-Pressure/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetStressAssessmentResultListHandle(c *gin.Context) {
	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize <= 0 {
		pageSize = 10
	}

	// 可选筛选
	var userID *int64
	if uidStr := c.Query("user_id"); uidStr != "" {
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err == nil && uid > 0 {
			userID = &uid
		}
	}

	var warningStatus *int
	if ws := c.Query("warning_status"); ws != "" {
		v, err := strconv.Atoi(ws)
		if err == nil {
			warningStatus = &v
		}
	}

	result, err := service.GetStressAssessmentList(service.StressQuery{
		UserID:        userID,
		StartDate:     c.Query("start_date"),
		EndDate:       c.Query("end_date"),
		WarningStatus: warningStatus,
		Page:          page,
		PageSize:      pageSize,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "获取成功",
		"data":    result,
	})
}
