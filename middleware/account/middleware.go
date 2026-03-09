package account

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func StateCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		//可以设置一些公共参数
		c.Set("example", "12345")
		//等其他中间件先执行
		c.Next()
		//获取耗时
		latency := time.Since(t)
		log.Printf("total cost time:%d us\n", latency/1000)
	}
}
