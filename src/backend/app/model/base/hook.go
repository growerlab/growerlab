package base

import (
	"fmt"
	"log"

	"github.com/growerlab/growerlab/src/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

type Action string
type Tense string

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

const (
	TenseBefore = "before"
	TenseAfter  = "after"
)

type Context struct {
	src    sqlx.Ext
	Table  string
	Values map[string]interface{} // key: column name, value: column value
}

type Act struct {
	Name              string               // 可以自定义一个名字
	Action            Action               // 订阅的动作
	Tense             Tense                // 动作的时态（before、after）
	Table             string               // 订阅的表
	Columns           []string             // 订阅的列
	CallbackFn        func(*Context) error // 执行回调
	IgnoreCallBackErr bool                 // 当执行回调失败，是否终止事务
}

var hook *Hook

func init() {
	hook = &Hook{
		hkSet: make(map[string]map[Action][]*Act),
	}
}

type Hook struct {
	hkSet map[string]map[Action][]*Act // table => action => []*Act
}

func (h *Hook) Register(act *Act) error {
	if act.CallbackFn == nil {
		return errors.New("hook callback can not be empty")
	}
	if act.Table == "" {
		return errors.New("hook table can not be empty")
	}
	if act.Action < ActionCreate || act.Action > ActionDelete {
		return errors.New("hook action can not be empty")
	}
	if act.Name == "" {
		act.Name = fmt.Sprintf("%s.%s.%s", act.Table, act.Action, act.Tense)
	}

	if _, found := h.hkSet[act.Table]; !found {
		h.hkSet[act.Table] = make(map[Action][]*Act)
	}
	h.hkSet[act.Table][act.Action] = append(h.hkSet[act.Table][act.Action], act)
	return nil
}

func (h *Hook) Effect(src sqlx.Ext, table string, a Action, t Tense, values map[string]interface{}) error {
	ctx := &Context{
		src:    src,
		Table:  table,
		Values: values,
	}

	if acts, found := h.hkSet[table]; found {
		if acts, found := acts[a]; found {
			for _, act := range acts {
				if act.CallbackFn == nil {
					continue
				}
				if act.Tense != t {
					continue
				}
				switch act.Action {
				case ActionCreate:
					err := act.CallbackFn(ctx)
					if err != nil {
						if !act.IgnoreCallBackErr {
							return err
						}
						log.Printf("[Hook] %+v\n", err)
					}
				case ActionUpdate:
					var match bool
					for _, col := range act.Columns {
						if _, match = values[col]; match {
							break
						}
					}
					if match {
						err := act.CallbackFn(ctx)
						if err != nil {
							if !act.IgnoreCallBackErr {
								return err
							}
							log.Printf("[Hook] %+v\n", err)
						}
					}
				case ActionDelete:
					err := act.CallbackFn(ctx)
					if err != nil {
						if !act.IgnoreCallBackErr {
							return err
						}
						log.Printf("[Hook] %+v\n", err)
					}
				}
			}
		}
	}
	return nil
}

func BeforeCreate(s *Act) error {
	s.Action = ActionCreate
	s.Tense = TenseBefore
	return hook.Register(s)
}

func AfterCreate(s *Act) error {
	s.Action = ActionCreate
	s.Tense = TenseAfter
	return hook.Register(s)
}

func BeforeUpdate(s *Act) error {
	s.Action = ActionUpdate
	s.Tense = TenseBefore
	return hook.Register(s)
}

func AfterUpdate(s *Act) error {
	s.Action = ActionUpdate
	s.Tense = TenseAfter
	return hook.Register(s)
}

func BeforeDelete(s *Act) error {
	s.Action = ActionDelete
	s.Tense = TenseBefore
	return hook.Register(s)
}

func AfterDelete(s *Act) error {
	s.Action = ActionDelete
	s.Tense = TenseAfter
	return hook.Register(s)
}
