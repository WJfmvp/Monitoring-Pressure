package admin

import (
	"net/http"
	"strconv"

	"Monitoring-Pressure/service"

	"github.com/gin-gonic/gin"
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

// PredictStressHandle 调用随机森林服务进行一次压力检测
// POST /admin/stress/predict
// body: StressPredictRequest 全部字段
func PredictStressHandle(c *gin.Context) {
	var req service.StressPredictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数格式错误: " + err.Error(),
		})
		return
	}

	resp, err := service.PredictAndSaveStress(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "压力预测失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "预测成功",
		"data": resp,
	})
}
