package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/service/common/session"
	"github.com/growerlab/growerlab/src/common/errors"
)

func Authorized(c *gin.Context) {
	sess := session.New(c)
	currentUser := sess.User()
	if currentUser == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.UnauthorizedError())
		return
	}
	c.Set("session", sess)
}
