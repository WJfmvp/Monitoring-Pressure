package account

import (
	"Monitoring-Pressure/util"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			util.ResponseError(c, util.ErrCodeNoPermission)
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			util.ResponseError(c, util.ErrCodeNoPermission)
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		util.ResponseError(c, util.ErrCodeNoPermission)
		c.Abort()
	}
}
