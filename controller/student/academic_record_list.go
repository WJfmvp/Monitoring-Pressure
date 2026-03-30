package student

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Monitoring-Pressure/service"
)

func GetMyAcademicRecordListHandle(c *gin.Context) {
	// 1. 从 JWT 中间件注入的上下文中获取当前登录用户ID
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未获取到当前登录用户信息",
		})
		return
	}

	var userID int64
	switch v := userIDValue.(type) {
	case int64:
		userID = v
	case int:
		userID = int64(v)
	case float64:
		userID = int64(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "当前登录用户信息格式错误",
		})
		return
	}

	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "当前登录用户ID不合法",
		})
		return
	}

	// 2. 读取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "page 参数不合法",
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "page_size 参数不合法",
		})
		return
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// 3. 可选筛选参数
	examName := c.Query("exam_name")
	term := c.Query("term")

	var minScorePtr *float64
	minScoreStr := c.Query("min_score")
	if minScoreStr != "" {
		minScore, err := strconv.ParseFloat(minScoreStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "min_score 参数不合法",
			})
			return
		}
		minScorePtr = &minScore
	}

	var maxScorePtr *float64
	maxScoreStr := c.Query("max_score")
	if maxScoreStr != "" {
		maxScore, err := strconv.ParseFloat(maxScoreStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "max_score 参数不合法",
			})
			return
		}
		maxScorePtr = &maxScore
	}

	if minScorePtr != nil && maxScorePtr != nil && *minScorePtr > *maxScorePtr {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "min_score 不能大于 max_score",
		})
		return
	}

	// 4. 调 service
	resp, err := service.GetStudentAcademicRecordList(service.StudentAcademicRecordListReq{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
		ExamName: examName,
		Term:     term,
		MinScore: minScorePtr,
		MaxScore: maxScorePtr,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询我的成绩列表失败: " + err.Error(),
		})
		return
	}

	// 5. 返回
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": resp,
	})
}
