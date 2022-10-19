package app

import (
	"sort"

	"github.com/growerlab/growerlab/src/common/errors"
)

var app *App

type Hook interface {
	Label() string                                               // 钩子名称
	Priority() uint                                              // 钩子优先级,数字越小越先执行
	Process(dispatcher EventDispatcher, sess *PushSession) error // 执行钩子
}

type App struct {
	dispatcher EventDispatcher
	hooks      []Hook
}

func (a *App) RegisterHook(hooks ...Hook) {
	a.hooks = append(a.hooks, hooks...)
	a.sortHooks()
}

func (a *App) sortHooks() {
	sort.Slice(a.hooks, func(i, j int) bool {
		return a.hooks[i].Priority() < a.hooks[j].Priority()
	})
}

func (a *App) Run(sess *PushSession) error {
	for _, hook := range a.hooks {
		if err := hook.Process(a.dispatcher, sess); err != nil {
			return err
		}
	}
	return nil
}

func Run(sess *PushSession) error {
	if app == nil {
		return errors.Errorf("must init App")
	}
	return app.Run(sess)
}
