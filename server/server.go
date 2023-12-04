package server

import (
	"godis/config"
	"godis/interface/redis"
	"godis/redis/parser"
	"strings"
	"sync/atomic"
)

type Server struct {
	dbSet []*atomic.Value // 数据库集合
}

func NewStandaloneServer() *Server {
	server := &Server{}
	// 默认创建16个数据库
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	server.dbSet = make([]*atomic.Value, config.Properties.Databases)
	for i := range server.dbSet {
		singleDB := makeDB()
		singleDB.index = i
		holder := &atomic.Value{}
		holder.Store(singleDB)
		server.dbSet[i] = holder
	}

	return server
}

func (s *Server) Exec(c redis.Connection, cmdLine [][]byte) redis.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))

	dbIndex := c.GetDBIndex()
	selectedDB, errReply := s.selectDB(dbIndex)
	if errReply != nil {
		return errReply
	}

	return selectedDB.Exec(c, cmdLine)
}

func (s *Server) selectDB(dbIndex int) (*DB, *parser.StandardErrReply) {
	if dbIndex >= len(s.dbSet) || dbIndex < 0 {
		return nil, parser.MakeErrReply("ERR DB index is out of range")
	}
	return s.dbSet[dbIndex].Load().(*DB), nil
}

func (s *Server) AfterClientClose(c redis.Connection) {

}

func (s *Server) Close() {

}
