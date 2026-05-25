package admin

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/users"
	"Monitoring-Pressure/util"

	"github.com/gin-gonic/gin"
)

// userListItem 不返回 password 字段
type userListItem struct {
	UserID    int64          `json:"user_id"`
	Username  string         `json:"username"`
	Telephone string         `json:"telephone"`
	Sex       int            `json:"sex"`
	Email     string         `json:"email"`
	Role      users.UserRole `json:"role"`
}

func GetUserListHandle(c *gin.Context) {
	var list []users.User
	err := db.DB.Find(&list).Error
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	resp := make([]userListItem, 0, len(list))
	for _, u := range list {
		resp = append(resp, userListItem{
			UserID:    u.UserID,
			Username:  u.Username,
			Telephone: u.Telephone,
			Sex:       u.Sex,
			Email:     u.Email,
			Role:      u.Role,
		})
	}

	util.ResponseSuccess(c, gin.H{
		"list": resp,
	})
}
