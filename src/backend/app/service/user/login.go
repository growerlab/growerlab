package user

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sessionModel "github.com/growerlab/growerlab/src/backend/app/model/session"
	userModel "github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/backend/app/utils/pwd"
	"github.com/growerlab/growerlab/src/backend/app/utils/uuid"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/asaskevich/govalidator.v9"
)

const TokenExpiredTime = 24 * time.Hour * 30 // 30天过期
const tokenField = "auth-user-token"

// Login 用户登录
//
//	 用户邮箱是否已验证
//		更新用户最后的登录时间/IP
//		生成用户登录token
func Login(ctx *gin.Context, req *LoginBasicAuth) (
	result *UserLoginResult,
	err error,
) {
	err = db.Transact(func(tx sqlx.Ext) error {
		loginService := NewLoginService(ctx.ClientIP(), db.DB)
		result, err = loginService.Do(req)
		if err != nil {
			return errors.Trace(err)
		}
		loginService.SetCookie(ctx)
		return nil
	})
	return result, err
}

type LoginBasicAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginService struct {
	ip string
	tx sqlx.Ext

	// session 登录完成后的session
	session *sessionModel.Session
}

func NewLoginService(ip string, tx sqlx.Ext) *LoginService {
	return &LoginService{
		ip: ip,
		tx: tx,
	}
}

func (l *LoginService) SetCookie(ctx *gin.Context) {
	ctx.SetCookie(tokenField, l.session.Token, 0, "/", "*", false, false)
}

func (r *LoginService) Do(auth *LoginBasicAuth) (
	result *UserLoginResult,
	err error,
) {
	user, err := r.Verify(auth)
	if err != nil {
		return nil, err
	}

	err = userModel.UpdateLogin(r.tx, user.ID, r.ip)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// 生成TOKEN返回给客户端
	r.session = r.buildAuthSession(user.ID, r.ip)
	err = sessionModel.New(r.tx).Add(r.session)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// namespace
	ns := user.Namespace()
	result = &UserLoginResult{
		Token:         r.session.Token,
		NamespacePath: ns.Path,
		Name:          user.Name,
		Email:         user.Email,
		PublicEmail:   user.PublicEmail,
	}
	return result, nil
}

func (r *LoginService) Verify(auth *LoginBasicAuth) (user *userModel.User, err error) {
	switch true {
	case !govalidator.IsByteLength(auth.Email, 1, 255):
		return nil, errors.InvalidParameterError(errors.User, errors.Email, errors.Empty)
	case !govalidator.IsByteLength(auth.Password, PasswordLenMin, PasswordLenMax):
		return nil, errors.InvalidParameterError(errors.User, errors.Password, errors.InvalidLength)
	}

	if strings.Contains(auth.Email, "@") {
		user, err = userModel.GetUserByEmail(r.tx, auth.Email)
		if err != nil {
			return nil, errors.Trace(err)
		}
	} else {
		user, err = userModel.GetUserByUsername(r.tx, auth.Email)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	if user == nil {
		return nil, errors.NotFoundError(errors.User)
	}
	if !user.Verified() {
		return nil, errors.AccessDenied(errors.User, errors.NotActivated)
	}

	ok := pwd.ComparePassword(user.EncryptedPassword, auth.Password)
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
