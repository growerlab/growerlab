package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/growerlab/growerlab/src/backend/app/model/activate"
	"github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
	"github.com/growerlab/growerlab/src/backend/app/utils/uuid"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/jmoiron/sqlx"
)

const ActivateExpiredTime = 24 * time.Hour

// 激活用户
func Activate(payload *ActivationCodePayload) (err error) {
	if !govalidator.IsByteLength(payload.Code, activate.CodeMaxLen, activate.CodeMaxLen) {
		return errors.P(errors.ActivationCode, errors.Code, errors.Invalid)
	}

	err = db.Transact(func(tx sqlx.Ext) error {
		err = DoActivate(tx, payload.Code)
		return err
	})
	return
}

// 激活账号的前期准备
// 生成code
// 生成url
// 生成模版
// 发送邮件
func DoPreActivate(tx sqlx.Ext, userID int64) error {
	code := buildActivateCode(userID)
	err := activate.AddCode(tx, code)
	if err != nil {
		return err
	}

	activateURL := buildActivateURL(code.Code)
	logger.Info("the activate url: %v", activateURL)

	// TODO 生成邮件模版(邮件模版功能应该抽出来独立，并能适配未来的其他模版)
	// TODO 发送邮件

	return nil
}

// 验证用户邮箱激活码
func DoActivate(tx sqlx.Ext, code string) error {
	acode, err := activate.GetCode(tx, code)
	if err != nil {
		return err
	}
	if acode == nil {
		return errors.NotFoundError(errors.ActivationCode)
	}
	// 是否已使用过
	if acode.UsedAt != nil {
		return errors.P(errors.ActivationCode, errors.Code, errors.Used)
	}
	// 是否过期
	// TODO 对于已经过期的激活码，应当在前端允许再次发送激活码（目前这块前后端还未开发）
	if acode.ExpiredAt < time.Now().Unix() {
		return errors.P(errors.ActivationCode, errors.Code, errors.Expired)
	}
	// 将code改成已使用
	err = activate.ActivateCode(tx, code)
	if err != nil {
		return err
	}
	// 激活用户状态
	err = user.ActivateUser(tx, acode.UserID)
	return err
}

func buildActivateURL(code string) string {
	baseURL := configurator.GetConf().WebsiteURL
	partURL := fmt.Sprintf("activate_user/%s", code)
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}
	return baseURL + partURL
}

func buildActivateCode(userID int64) *activate.ActivationCode {
	code := new(activate.ActivationCode)
	code.UserID = userID
	code.Code = uuid.UUIDv16()
	code.ExpiredAt = time.Now().Add(ActivateExpiredTime).Unix()
	return code
}
