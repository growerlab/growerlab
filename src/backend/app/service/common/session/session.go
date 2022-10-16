package session

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/common/env"
	userModel "github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
	"github.com/growerlab/growerlab/src/common/db"
)

const (
	AuthUserToken = "auth-user-token"
)

type Session struct {
	environment *env.Environment
	ctx         *gin.Context
	user        *userModel.User
}

func New(c *gin.Context) *Session {
	var e = env.NewEnvironment()
	var user *userModel.User
	var userToken = GetUserToken(c)
	var err error

	if len(userToken) > 0 {
		user, err = userModel.GetUserByUserToken(db.DB, userToken)
		if err != nil {
			logger.Error("get user by user token failed, user token: %s, err: %s", userToken, err.Error())
			return nil
		}
	}

	e.Set(env.VarUserToken, userToken)
	return &Session{
		environment: e,
		ctx:         c,
		user:        user,
	}
}

func (s *Session) GetContext() *gin.Context {
	return s.ctx
}

func (s *Session) Env() *env.Environment {
	return s.environment
}

func (s *Session) User() *userModel.User {
	return s.user
}

func (s *Session) UserNamespace() *int64 {
	if s.user == nil {
		return nil
	}
	return &s.user.NamespaceID
}

func (s *Session) Token() string { // current user
	userToken, _ := s.environment.MustString(env.VarUserToken)
	return userToken
}

func (s *Session) IsGuest() bool {
	token := s.Token()
	return len(token) == 0
}

func GetUserToken(ctx *gin.Context) string {
	token := getValueFromHeaderOrCookie(AuthUserToken, ctx)
	return token
}

func getValueFromHeaderOrCookie(k string, ctx *gin.Context) string {
	v := ctx.GetHeader(k)
	if len(v) < 5 {
		v, _ = ctx.Cookie(k)
	}
	return v
}
