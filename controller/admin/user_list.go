package admin

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/users"
	"Monitoring-Pressure/util"

	"github.com/gin-gonic/gin"
)

func GetUserListHandle(c *gin.Context) {
	var list []users.UserInfo
	err := db.DB.Find(&list).Error
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"list": list,
	})
}
