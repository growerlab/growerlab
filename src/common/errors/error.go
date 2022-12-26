package errors

import (
	"fmt"
	"net/http"
	"strings"

	pkgerr "github.com/pkg/errors"
)

// TODO 目前有一些错误不应该在后端输出错误，例如 NotFound， 此类错误是没有必要的

// 定义错误
const (
	// 非法参数
	invalidParameter = "InvalidParameter"
	// 缺少参数
	missingParameter = "MissingParameter"
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
	unauthorizedError = "UnauthorizedError"
	// PermissionError
	permissionError = "PermissionError"
	// 仓库
	repositoryError = "RepositoryError"
	//
	timeoutError = "TimeoutError"
)

// 定义错误原因
const (
	// Invalid 非法的
	Invalid = "Invalid"
	// NotFoundField 无法找到属性（字段）
	NotFoundField = "NotFoundField"
	// InvalidLength 非法长度
	InvalidLength = "InvalidLength"
	// Expired 失效，过期
	Expired = "Expired"
	// Used 已被使用过
	Used = "Used"
	// NotEqual 不匹配
	NotEqual = "NotEqual"
	// Empty 空的
	Empty = "Empty"
	// AlreadyExists 已存在
	AlreadyExists = "AlreadyExists"
	// NotActivated 未激活
	NotActivated = "NotActivated"
	// NoPermission 无权限
	NoPermission = "NoPermission"
	// GitServerInvalid git服务器异常
	GitServerInvalid = "GitServerInvalid"
)

var httpCodeSet = map[string]int{
	invalidParameter:  http.StatusBadRequest,
	missingParameter:  http.StatusBadRequest,
	notFoundError:     http.StatusNotFound,
	graphQLError:      http.StatusBadRequest,
	alreadyExists:     http.StatusConflict,
	accessDeniedError: http.StatusForbidden,
	sqlError:          http.StatusInternalServerError,
	unauthorizedError: http.StatusUnauthorized,
	permissionError:   http.StatusForbidden,
	repositoryError:   http.StatusBadRequest,
	timeoutError:      http.StatusRequestTimeout,
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

func MissingParameterError(model, field string) error {
	return mustCode(nil, missingParameter, field)
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

func UnauthorizedError() error {
	return mustCode(nil, unauthorizedError)
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

func TimeoutError(reason string) error {
	return mustCode(nil, timeoutError, reason)
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
	Wrap = func(err error, message string) error {
		err = pkgerr.Cause(err)
		return pkgerr.Wrap(err, message)
	}
	Wrapf = func(err error, format string, args ...interface{}) error {
		err = pkgerr.Cause(err)
		return pkgerr.Wrapf(err, format, args...)
	}
	Message = func(err error, message string) error {
		err = pkgerr.Cause(err)
		return pkgerr.WithMessage(err, message)
	}
	Messagef = func(err error, format string, args ...interface{}) error {
		err = pkgerr.Cause(err)
		return pkgerr.WithMessagef(err, format, args...)
	}
	Trace = func(err error) error {
		err = pkgerr.Cause(err)
		return pkgerr.WithStack(err)
	}
	WithStack = Trace
	Cause     = pkgerr.Cause
	Errorf    = pkgerr.Errorf
	New       = pkgerr.New
	Is        = pkgerr.Is
)
