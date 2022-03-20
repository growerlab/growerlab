package errors

import (
	"fmt"
	"strings"

	pkgerr "github.com/pkg/errors"
)

// TODO 目前有一些错误不应该在后端输出错误，例如 NotFound， 此类错误是没有必要的

// 定义错误
const (
	// 非法参数
	invalidParameter = "InvalidParameter"
	// 无法找到
	notFoundError = "NotFoundError"
	// GraphQLError
	graphQLError = "GraphQLError"
	// 已存在
	alreadyExists = "AlreadyExists"
	// AccessDenied
	accessDeniedError = "AccessDeniedError"
	// sql错误
	sqlError = "SQLError"
	// 未登录
	unauthorized = "Unauthorized"
	// PermissionError
	permissionError = "PermissionError"
	// 仓库
	repositoryError = "RepositoryError"
)

// 定义错误原因
const (
	// 非法的
	Invalid = "Invalid"
	// 无法找到属性（字段）
	NotFoundField = "NotFoundField"
	// 非法长度
	InvalidLength = "InvalidLength"
	// 失效，过期
	Expired = "Expired"
	// 已被使用过
	Used = "Used"
	// 不匹配
	NotEqual = "NotEqual"
	// 空的
	Empty = "Empty"
	// 已存在
	AlreadyExists = "AlreadyExists"
	// 未激活
	NotActivated = "NotActivated"
	// 仓库服务异常
	SvcServerNotReady = "SvcServerNotReady"
	// 无权限
	NoPermission = "NoPermission"
)

var httpCodeSet = map[string]int{
	invalidParameter:  400,
	notFoundError:     404,
	graphQLError:      400,
	alreadyExists:     409,
	accessDeniedError: 403,
	sqlError:          500,
	unauthorized:      401,
	permissionError:   403,
	repositoryError:   500,
}

type Result struct {
	Err        error  `json:"-"`
	Code       string `json:"code"`
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e *Result) Error() string {
	var sb = strings.Builder{}
	sb.WriteString(e.Message)
	if e.Err != nil {
		sb.WriteString(": ")
		sb.WriteString(e.Err.Error())
	}
	return sb.String()
}

var P = InvalidParameterError

func InvalidParameterError(model, field, reason string) error {
	return mustCode(nil, invalidParameter, model, field, reason)
}

func NotFoundError(model string) error {
	return mustCode(nil, notFoundError, model)
}

func AlreadyExistsError(model, reason string) error {
	return mustCode(nil, alreadyExists, model, reason)
}

func SQLError(err error) error {
	return mustErr(err, sqlError)
}

func GraphQLError() error {
	return mustCode(nil, graphQLError)
}

func Unauthorize() error {
	return mustCode(nil, unauthorized)
}

func AccessDenied(model, reason string) error {
	return mustCode(nil, accessDeniedError, model, reason)
}

func PermissionError(reason string) error {
	return mustCode(nil, permissionError, reason)
}

func RepositoryError(reason string) error {
	return mustCode(nil, repositoryError, reason)
}

func mustErr(err error, parts ...string) error {
	if err == nil {
		return nil
	}
	return mustCode(err, parts...)
}

// 必须调用该方法生成<xxx>字符串，便于前端解析数据
func mustCode(err error, parts ...string) error {
	if len(parts) == 0 {
		panic("parts is required")
	}
	hc, ok := httpCodeSet[parts[0]]
	if !ok {
		hc = 500
	}
	return Trace(&Result{
		Err:        err,
		Code:       parts[0],
		StatusCode: hc,
		Message:    fmt.Sprintf("<%s>", strings.Join(parts, ".")),
	})
}

// 封装（避免在项目中使用时，引用多个包）
var (
	Wrap     = pkgerr.Wrap
	Wrapf    = pkgerr.Wrapf
	Message  = pkgerr.WithMessage
	Messagef = pkgerr.WithMessagef
	Trace    = pkgerr.WithStack
	Cause    = pkgerr.Cause
	Errorf   = pkgerr.Errorf
	New      = pkgerr.New
)
