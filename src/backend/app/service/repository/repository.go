package repository

import (
	"github.com/gin-gonic/gin"
	userModel "github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/backend/app/service/common/session"
)

type Repository struct {
	currentUser   *userModel.User
	currentUserID *int64
	namespace     string
	// 当取list时，repo可以为空
	repo *string
}

func New(c *gin.Context, namespace string, path *string) *Repository {
	u := session.New(c).User()
	currentUserID := session.New(c).UserID()
	return &Repository{
		currentUser:   u,
		currentUserID: currentUserID,
		namespace:     namespace,
		repo:          path,
	}
}
