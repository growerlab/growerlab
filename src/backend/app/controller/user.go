package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/growerlab/growerlab/src/backend/app/service/user"
)

func RegisterUser(c *gin.Context) {
	var req user.NewUserPayload
	if err := c.BindJSON(&req); err != nil {
		Render(c, nil, err)
		return
	}

	clientIP := c.ClientIP()
	err := user.Register(&req, clientIP)
	Render(c, nil, err)
}

func ActivateUser(c *gin.Context) {
	var req user.ActivationCodePayload
	if err := c.BindJSON(&req); err != nil {
		Render(c, nil, err)
		return
	}
	err := user.Activate(&req)
	Render(c, nil, err)
}

func LoginUser(c *gin.Context) {
	var input user.LoginBasicAuth
	if err := c.BindJSON(&input); err != nil {
		Render(c, nil, err)
		return
	}
	result, err := user.Login(c, &input)
	Render(c, result, err)
}
