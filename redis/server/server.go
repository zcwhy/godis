package server

import (
	"context"
	"godis/lib/sync/atomic"
	"net"
	"sync"
)

type Server struct {
	closing    atomic.Boolean
	activeConn sync.Map
}

func MakeServer() *Server {
	return &Server{}
}

func (h *Server) Handle(ctx context.Context, conn net.Conn) {
	client := connection.NewConn(conn)
}
