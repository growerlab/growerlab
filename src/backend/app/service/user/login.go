package user

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
	"github.com/growerlab/growerlab/src/backend/app/model/db"
	sessionModel "github.com/growerlab/growerlab/src/backend/app/model/session"
	userModel "github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/backend/app/utils/pwd"
	"github.com/growerlab/growerlab/src/backend/app/utils/uuid"
	"github.com/jmoiron/sqlx"
	"gopkg.in/asaskevich/govalidator.v9"
)

const TokenExpiredTime = 24 * time.Hour * 30 // 30天过期
const tokenField = "auth-user-token"

// Login 用户登录
//  用户邮箱是否已验证
//	更新用户最后的登录时间/IP
//	生成用户登录token
func Login(ctx *gin.Context, req *LoginBasicAuth) (
	result *UserLoginResult,
	err error,
) {
	loginService := NewLoginService(ctx.ClientIP(), req)
	result, err = loginService.Do(db.DB)
	if err != nil {
		return nil, err
	}
	loginService.SetCookie(ctx)
	return
}

type LoginBasicAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginService struct {
	ip   string
	auth *LoginBasicAuth

	// session 登录完成后的session
	session *sessionModel.Session
}

func NewLoginService(ip string, auth *LoginBasicAuth) *LoginService {
	return &LoginService{
		ip:   ip,
		auth: auth,
	}
}

func (l *LoginService) SetCookie(ctx *gin.Context) {
	ctx.SetCookie(tokenField, l.session.Token, 0, "/", ctx.Request.Host, false, false)
}

func (l *LoginService) Do(src sqlx.Ext) (
	result *UserLoginResult,
	err error,
) {
	user, err := l.prepare(src)
	if err != nil {
		return nil, err
	}

	err = db.Transact(func(tx sqlx.Ext) error {
		err = userModel.UpdateLogin(tx, user.ID, l.ip)
		if err != nil {
			return err
		}

		// 生成TOKEN返回给客户端
		l.session = l.buildAuthSession(user.ID, l.ip)
		err = sessionModel.New(tx).Add(l.session)
		if err != nil {
			return err
		}

		// namespace
		ns := user.Namespace()
		result = &UserLoginResult{
			Token:         l.session.Token,
			NamespacePath: ns.Path,
			Name:          user.Name,
			Email:         user.Email,
			PublicEmail:   user.PublicEmail,
		}
		return nil
	})
	return result, err
}

func (r *LoginService) prepare(src sqlx.Queryer) (user *userModel.User, err error) {
	switch true {
	case !govalidator.IsByteLength(r.auth.Email, 1, 255):
		return nil, errors.InvalidParameterError(errors.User, errors.Email, errors.Empty)
	case !govalidator.IsByteLength(r.auth.Password, PasswordLenMin, PasswordLenMax):
		return nil, errors.InvalidParameterError(errors.User, errors.Password, errors.InvalidLength)
	}

	if strings.Contains(r.auth.Email, "@") {
		user, err = userModel.GetUserByEmail(src, r.auth.Email)
		if err != nil {
			return nil, err
		}
	} else {
		user, err = userModel.GetUserByUsername(src, r.auth.Email)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, errors.NotFoundError(errors.User)
	}
	if !user.Verified() {
		return nil, errors.AccessDenied(errors.User, errors.NotActivated)
	}

	ok := pwd.ComparePassword(user.EncryptedPassword, r.auth.Password)
	if !ok {
		return nil, errors.InvalidParameterError(errors.User, errors.Password, errors.NotEqual)
	}
	return user, nil
}

func (r *LoginService) buildAuthSession(userID int64, clientIP string) *sessionModel.Session {
	return &sessionModel.Session{
		OwnerID:   userID,
		Token:     uuid.UUID(),
		ClientIP:  clientIP,
		CreatedAt: time.Now().Unix(),
		ExpiredAt: time.Now().Add(TokenExpiredTime).Unix(),
	}
}
