package env

import (
	"github.com/growerlab/growerlab/src/common/errors"
	"strconv"
)

const (
	VarUserToken = "user-token"
)

var (
	ErrNotExists = errors.New("not exists")
)

type Environment struct {
	val map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		val: make(map[string]interface{}),
	}
}

func (e *Environment) Set(k string, v interface{}) {
	e.val[k] = v
}

func (e *Environment) Get(k string) interface{} {
	return e.val[k]
}

func (e *Environment) String(k string) (string, bool) {
	v, ok := e.Get(k).(string)
	return v, ok
}

func (e *Environment) MustString(k string) (v string, err error) {
	v, ok := e.String(k)
	if ok {
		return v, nil
	}
	return "", errors.Trace(ErrNotExists)
}

func (e *Environment) Int64(k string) (int64, bool) {
	v := e.Get(k)
	var i int64
	switch v.(type) {
	case int64:
		i = v.(int64)
	case int:
		i = int64(v.(int))
	case string:
		var err error
		i, err = strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return 0, false
		}
	default:
		return 0, false
	}
	return i, true
}

func (e *Environment) MustInt64(k string) (v int64, err error) {
	v, ok := e.Int64(k)
	if ok {
		return v, nil
	}
	return 0, errors.Trace(ErrNotExists)
}
