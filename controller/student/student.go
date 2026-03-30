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
	// 获取用户ID，通常是通过 JWT 解码获取
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		log.Printf("无效请求：未提供用户ID，无法获取干预建议记录")
		return
	}

	// 调用服务层获取该用户的所有干预建议记录
	interventions, err := service.GetInterventionRecordsByUserID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取干预建议记录失败: %v", err),
		})
		log.Printf("获取用户 %d 干预建议记录失败: %v", userID, err)
		return
	}

	// 如果没有找到记录，返回空响应或提示信息
	if len(interventions) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有找到干预建议记录",
		})
		return
	}

	// 返回查询到的数据
	c.JSON(http.StatusOK, gin.H{
		"message": "获取干预建议记录成功",
		"data":    interventions,
	})
}
