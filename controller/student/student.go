package student

import (
	"Monitoring-Pressure/util"
	"github.com/gin-gonic/gin"
)

func GetHomeHandle(c *gin.Context) {
	util.ResponseSuccess(c, gin.H{
		"message": "学生首页接口",
	})
}
