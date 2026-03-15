package account

import (
	myjwt "Monitoring-Pressure/JWT"
	"Monitoring-Pressure/util"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.ResponseError(c, util.ErrCodeNeedLogin)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.ResponseError(c, util.ErrCodeNeedLogin)
			c.Abort()
			return
		}

		claims, err := myjwt.ParseToken(parts[1])
		if err != nil {
			util.ResponseError(c, util.ErrCodeNeedLogin)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("telephone", claims.Telephone)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}
