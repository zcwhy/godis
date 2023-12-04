package net

import (
	"context"
	"godis/interface/redis"
	"godis/lib/log"
	"godis/lib/sync/atomic"
	"godis/redis/connection"
	"godis/redis/parser"
	"godis/server"
	"net"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type Handler struct {
	activeConn sync.Map
	closing    atomic.Boolean
	server     redis.Server
}

func MakeHandler() *Handler {
	return &Handler{
		server: server.NewStandaloneServer(),
	}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	client := connection.NewConn(conn)
	h.activeConn.Store(client, struct{}{})

	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {

		}
		if payload.Data == nil {
			log.Error("empty payload")
			continue
		}

		// 从客户端接收到multiBulkReply类型的命令
		r, ok := payload.Data.(*parser.MultiBulkReply)
		if !ok {
			log.Error("require multi bulk protocol")
			continue
		}
		result := h.server.Exec(client, r.Args)
		if result != nil {
			_, _ = client.Write(result.ToBytes())
		} else {
			_, _ = client.Write(unknownErrReplyBytes)
		}
	}

}
