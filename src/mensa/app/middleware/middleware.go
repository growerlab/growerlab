package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/growerlab/growerlab/src/mensa/app/common"
)

type MiddlewareError string

func (m MiddlewareError) Error() string {
	return string(m)
}

type HandleResult struct {
	status     int // http status
	lastError  error
	gitMessage strings.Builder
}

func (h *HandleResult) Status() int {
	return h.status
}

// 当进入失败时，应返回http错误的信息
func (h *HandleResult) Message() string {
	if h.status > http.StatusCreated {
		h.gitMessage.WriteString(fmt.Sprintf("\n%d", h.status))
	}
	if h.lastError != nil {
		h.gitMessage.WriteString(h.lastError.Error())
	}
	return h.gitMessage.String()
}

// 错误码
func (h *HandleResult) LastErr() error {
	return h.lastError
}

type HandleFunc func(*common.Context) (httpCode int, appendText string, err error)

type Middleware struct {
	funcs []HandleFunc
}

func (m *Middleware) Add(fn HandleFunc) {
	m.funcs = append(m.funcs, fn)
}

func (m *Middleware) run(ctx *common.Context) *HandleResult {
	result := &HandleResult{}

	for _, fn := range m.funcs {
		statusCode, appendText, err := fn(ctx)
		if len(appendText) > 0 {
			result.gitMessage.WriteString(appendText)
		}
		result.status = statusCode
		if err != nil {
			result.lastError = err
			return result
		}
	}
	return result
}

func (m *Middleware) Enter(ctx *common.Context) *HandleResult {
	return m.run(ctx)
}
