package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
	"github.com/growerlab/growerlab/src/common/errors"
)

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
