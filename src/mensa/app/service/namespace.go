package service

import (
	"strconv"

	userModel "github.com/growerlab/growerlab/src/backend/app/model/user"
	userService "github.com/growerlab/growerlab/src/backend/app/service/user"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/mensa/app/common"
)

func GetNamespaceByOperator(operator *common.Operator) (int64, error) {
	if operator.IsHttp() {
		username := operator.HttpUser.Username()
		password, pwdExists := operator.HttpUser.Password()
		if !pwdExists {
			return 0, errors.New("password is required")
		}
		user, err := userService.NewLoginService("", db.DB).Verify(&userService.LoginBasicAuth{
			Email:    username,
			Password: password,
		})
		if err != nil {
			return 0, errors.Trace(err)
		}
		return user.NamespaceID, nil
	} else { // ssh
		// TODO SSH
		return 0, errors.New("ssh ...")
	}
}

func GetUserNamespaceByUsername(username string) (int64, error) {
	key := db.MemDB.KeyMaker().Append("user", "namespace").String()
	field := username

	userNamespaceID, err := NewCache().GetOrSet(
		key,
		field,
		func() (string, error) {
			u, err := userModel.GetUserByUsername(db.DB, username)
			if err != nil {
				return "", err
			}
			if u == nil {
				return "", errors.Errorf("not found user: %s", username)
			}
			return strconv.FormatInt(u.NamespaceID, 10), nil
		})
	if err != nil {
		return 0, err
	}
	if userNamespaceID == "" || userNamespaceID == "0" {
		return 0, errors.Errorf("not found user: %s.", username)
	}
	return strconv.ParseInt(userNamespaceID, 10, 64)
}
