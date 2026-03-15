package admin

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/users"
	"Monitoring-Pressure/util"

	"github.com/gin-gonic/gin"
)

type UpdateUserRoleRequest struct {
	Telephone string         `json:"telephone"`
	Role      users.UserRole `json:"role"`
}

func UpdateUserRoleHandle(c *gin.Context) {
	var req UpdateUserRoleRequest
	if err := c.BindJSON(&req); err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if req.Telephone == "" || req.Role == "" {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if req.Role != users.RoleStudent &&
		req.Role != users.RoleTeacher &&
		req.Role != users.RoleAdmin {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	result := db.DB.Model(&users.User{}).
		Where("telephone = ?", req.Telephone).
		Update("role", req.Role)

	if result.Error != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	if result.RowsAffected == 0 {
		util.ResponseError(c, util.ErrCodeUserNotExist)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"message":   "角色修改成功",
		"telephone": req.Telephone,
		"role":      req.Role,
	})
}
