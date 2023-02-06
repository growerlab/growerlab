package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/service/repository"
)

func Repositories(c *gin.Context) {
	namespace := c.Param("namespace")
	repos, err := repository.New(c, namespace, nil).List()
	Render(c, repos, err)
}

func Repository(c *gin.Context) {
	namespace := c.Param("namespace")
	repo := c.Param("repo")

	r, err := repository.New(c, namespace, &repo).Get()
	Render(c, r, err)
}

func RepositoryTree(c *gin.Context) {
	namespace := c.Param("namespace")
	repo := c.Param("repo")
	ref := c.Param("ref")
	folder := c.Param("folder")
	r, err := repository.New(c, namespace, &repo).TreeFiles(ref, &folder)
	Render(c, r, err)
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
