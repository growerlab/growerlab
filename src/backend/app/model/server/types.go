package server

import "fmt"

// 部署svc的服务器

type Server struct {
	ID        int64  `db:"id"`
	Summary   string `db:"summary"`
	Host      string `db:"host"`
	Port      int    `db:"port"`
	State     int    `db:"state"`
	CreatedAt int64  `db:"created_at"`
	DeletedAt *int64 `db:"deleted_at"`
}

func (s *Server) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func (s *Server) ApiURL() string {
	return fmt.Sprintf("%s/graphql", s.URL())
}
