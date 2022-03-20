package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/service/repository"
)

func Repositories(c *gin.Context) {
	namespace := c.Param("namespace")
	repos, err := repository.ListRepositories(c, namespace)
	Render(c, repos, err)
}

func Repository(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	repo, err := repository.GetRepository(c, namespace, name)
	Render(c, repo, err)
}

func CreateRepository(c *gin.Context) {
	var req repository.NewRepositoryPayload
	if err := c.BindJSON(&req); err != nil {
		Render(c, nil, err)
		return
	}

	err := repository.CreateRepository(c, &req)
	Render(c, nil, err)
}
