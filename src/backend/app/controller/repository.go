package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/service/repository"
)

func Repositories(c *gin.Context) {
	namespace := c.Param("namespace")
	repos, err := repository.NewTaker(c, namespace, nil).List()
	Render(c, repos, err)
}

func Repository(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	repo, err := repository.NewTaker(c, namespace, &name).Get()
	Render(c, repo, err)
}

func CreateRepository(c *gin.Context) {
	var req repository.CreateParams
	if err := c.BindJSON(&req); err != nil {
		Render(c, nil, err)
		return
	}

	err := repository.NewCreator(c, &req).Do()
	Render(c, nil, err)
}
