package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSForLocal 处理本地访问的CORS
func CORSForLocal(c *gin.Context) {
	// if !conf.GetConf().Debug {
	// 	return
	// }
	reqAccessHeaders := c.Request.Header.Get("Access-Control-Request-Headers")

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", reqAccessHeaders)
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
