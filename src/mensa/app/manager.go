package app

import (
	"runtime/debug"
	"sync"

	"github.com/growerlab/growerlab/src/mensa/app/common"
)

type MiddlewareResult struct {
	HttpCode    int
	HttpMessage string
	Err         error
}

type MiddlewareHandler func(ctx *common.Context) *MiddlewareResult

type Server interface {
	// 启动并监听服务
	// 	当有新的链接时，将调用cb方法
	ListenAndServe(MiddlewareHandler) error
	// 停止服务
	Shutdown() error
}

type Manager struct {
	servers []Server

	entry Entryer
}

func NewManager(entry Entryer) *Manager {
	return &Manager{
		entry: entry,
	}
}

func (m *Manager) RegisterServer(srv Server) {
	m.servers = append(m.servers, srv)
}

// run server and waiting for end
func (m *Manager) Run() {
	var wg sync.WaitGroup
	for _, s := range m.servers {
		wg.Add(1)
		go func(srv Server) {
			defer wg.Done()
			defer func() {
				if e := recover(); e != nil {
					debug.PrintStack()
				}
			}()
			err := srv.ListenAndServe(m.ServerHandler)
			if err != nil {
				panic(err)
			}
		}(s)
	}
	wg.Wait()
}

func (m *Manager) ServerHandler(ctx *common.Context) *MiddlewareResult {
	if m.entry != nil {
		result := m.entry.Enter(ctx)
		return &MiddlewareResult{
			HttpCode:    result.Status(),
			HttpMessage: result.Message(),
			Err:         result.LastErr(),
		}
	}
	return nil
}
