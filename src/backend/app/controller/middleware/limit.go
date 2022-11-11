package middleware

import (
	"net/http"

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
