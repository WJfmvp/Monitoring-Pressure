package admin

import (
	"Monitoring-Pressure/util"

	"github.com/gin-gonic/gin"
)

func GetHomeHandle(c *gin.Context) {
	util.ResponseSuccess(c, gin.H{
		"message": "管理员首页接口",
	})
}
