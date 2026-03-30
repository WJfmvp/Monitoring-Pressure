package student

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Monitoring-Pressure/service"
)

func GetLatestStressAssessmentResultHandle(c *gin.Context) {
	// 1. 获取当前登录用户ID
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
			"msg":  "用户ID类型错误",
		})
		return
	}

	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户ID不合法",
		})
		return
	}

	// 2. 调 service
	resp, err := service.GetStudentLatestStressResult(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 3. 没数据
	if resp == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "暂无压力评估结果",
			"data": nil,
		})
		return
	}

	// 4. 返回
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": resp,
	})
}
