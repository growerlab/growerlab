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
	repo := c.Param("repo")

	r, err := repository.NewTaker(c, namespace, &repo).Get()
	Render(c, r, err)
}

func RepositoryTree(c *gin.Context) {
	namespace := c.Param("namespace")
	repo := c.Param("repo")
	ref := c.Param("ref")
	dir := c.Param("dir")
	r, err := repository.NewTaker(c, namespace, &repo).TreeFiles(ref, &dir)
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
