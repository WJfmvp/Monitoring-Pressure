package teacher

import (
	"Monitoring-Pressure/util"
	"github.com/gin-gonic/gin"
)

func GetHomeHandle(c *gin.Context) {
	util.ResponseSuccess(c, gin.H{
		"message": "老师首页接口",
	})
}
