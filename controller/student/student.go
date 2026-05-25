package student

import (
	"Monitoring-Pressure/service"
	"Monitoring-Pressure/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetHomeHandle(c *gin.Context) {
	util.ResponseSuccess(c, gin.H{
		"message": "学生首页接口",
	})
}

// GetMyInterventionRecordListHandle 获取我的干预建议记录列表
func GetMyInterventionRecordListHandle(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}

	interventions, err := service.GetInterventionRecordsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取干预建议记录失败: %v", err),
		})
		log.Printf("获取用户 %d 干预建议记录失败: %v", userID, err)
		return
	}

	// 没记录也走统一返回结构，避免前端字段不一致
	c.JSON(http.StatusOK, gin.H{
		"message": "获取干预建议记录成功",
		"data":    interventions,
	})
}

// getUserIDFromContext 安全地从 gin.Context 拿到 int64 user_id
func getUserIDFromContext(c *gin.Context) (int64, bool) {
	v, exists := c.Get("user_id")
	if !exists || v == nil {
		return 0, false
	}
	switch id := v.(type) {
	case int64:
		return id, id > 0
	case int:
		return int64(id), id > 0
	case float64:
		return int64(id), id > 0
	default:
		return 0, false
	}
}
