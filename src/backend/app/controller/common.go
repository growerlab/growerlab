package controller

import (
	"github.com/growerlab/growerlab/src/common/errors"
	"net/http"

	"github.com/growerlab/growerlab/src/backend/app/utils/logger"

	"github.com/gin-gonic/gin"
)

const (
	MaxGraphQLRequestBody = int64(1 << 20) // 1MB
)

func LimitGETRequestBody(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		return
	}
	if ctx.Request.ContentLength > MaxGraphQLRequestBody {
		ctx.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
}

func Render(c *gin.Context, payload interface{}, err error) {
	if err != nil {
		cerr := errors.Cause(err)
		if e, ok := cerr.(*errors.Result); ok {
			c.AbortWithStatusJSON(e.StatusCode, cerr)

			if e2 := errors.Cause(e.Err); e2 != nil {
				logger.Error("render2: %+v\n", e2)
			}
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, cerr)
		}
		logger.Error("render: %+v\n", err)

		return
	}
	if payload != nil {
		c.AbortWithStatusJSON(http.StatusOK, payload)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, &errors.Result{
		Code: "ok",
	})
}

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
