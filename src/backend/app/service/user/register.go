package user

import (
	"time"

	nsModel "github.com/growerlab/growerlab/src/backend/app/model/namespace"
	userModel "github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/backend/app/utils/pwd"
	"github.com/growerlab/growerlab/src/backend/app/utils/regex"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/asaskevich/govalidator.v9"
)

const (
	PasswordLenMin = 7
	PasswordLenMax = 32

	UsernameLenMin = 4
	UsernameLenMax = 40
)

type ActivationCodePayload struct {
	Code string `json:"code"`
}

type NewUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserLoginResult struct {
	Token       string `json:"token"`
	Namespace   string `json:"namespace"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	PublicEmail string `json:"public_email"`
}

func validateRegisterUser(payload *NewUserPayload) error {
	if !govalidator.IsEmail(payload.Email) {
		return errors.P(errors.User, errors.Email, errors.Invalid)
	}
	if !govalidator.IsByteLength(payload.Password, PasswordLenMin, PasswordLenMax) {
		return errors.P(errors.User, errors.Password, errors.InvalidLength)
	}
	if !govalidator.IsByteLength(payload.Username, UsernameLenMin, UsernameLenMax) {
		return errors.P(errors.User, errors.Username, errors.InvalidLength)
	}
	if !regex.Match(payload.Username, regex.UsernameRegex) {
		return errors.P(errors.User, errors.Username, errors.Invalid)
	}
	if !regex.Match(payload.Password, regex.PasswordRegex) {
		return errors.P(errors.User, errors.Password, errors.Invalid)
	}

	// 不允许使用的关键字
	if _, invalidUsername := userModel.InvalidUsernameSet[payload.Username]; invalidUsername {
		return errors.AlreadyExistsError(errors.User, errors.AlreadyExists)
	}

	// email, username是否已经存在
	exists, err := userModel.ExistsEmailOrUsername(db.DB, payload.Username, payload.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.AlreadyExistsError(errors.User, errors.AlreadyExists)
	}
	return nil
}

func buildUser(payload *NewUserPayload, clientIP string) (*userModel.User, error) {
	password, err := pwd.GeneratePassword(payload.Password)
	if err != nil {
		return nil, err
	}
	return &userModel.User{
		Email:             payload.Email,
		EncryptedPassword: password,
		Username:          payload.Username,
		Name:              payload.Username,
		PublicEmail:       payload.Email,
		CreatedAt:         time.Now().Unix(),
		RegisterIP:        clientIP,
		IsAdmin:           false,
		NamespaceID:       0,
	}, nil
}

func buildNamespace(user *userModel.User) *nsModel.Namespace {
	return &nsModel.Namespace{
		Path:    user.Username,
		OwnerID: user.ID,
		Type:    int(nsModel.TypeUser),
	}
}

// Register 用户注册
// 1. 将用户信息添加到数据库中
// 2. 发送验证邮件（这里可以考虑使用KeyDB来建立邮件发送队列，避免重启进程后，发送任务丢失）
// 3. Done
func Register(payload *NewUserPayload, clientIP string) error {
	var err error
	err = validateRegisterUser(payload)
	if err != nil {
		return err
	}

	err = db.Transact(func(tx sqlx.Ext) error {
		user, err := buildUser(payload, clientIP)
		if err != nil {
			return err
		}

		err = userModel.AddUser(tx, user)
		if err != nil {
			return err
		}

		// create namespace
		ns := buildNamespace(user)
		err = nsModel.AddNamespace(tx, ns)
		if err != nil {
			return err
		}

		// set namespace id to user
		err = userModel.UpdateNamespace(tx, user.ID, ns.ID)
		if err != nil {
			return err
		}

		// activate user
		err = DoPreActivate(tx, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
