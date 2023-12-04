package connection

import (
	"godis/lib/log"
	"godis/lib/sync/wait"
	"net"
	"sync"
)

type Connection struct {
	conn        net.Conn
	selectedDB  int // 连接的数据库
	sendingData wait.Wait
}

var connPool = sync.Pool{
	New: func() interface{} {
		return &Connection{}
	},
}

// NewConn creates Connection instance
func NewConn(conn net.Conn) *Connection {
	c, ok := connPool.Get().(*Connection)
	if !ok {
		log.Error("connection pool make wrong type")
		return &Connection{
			conn: conn,
		}
	}
	c.conn = conn
	return c
}

func (c *Connection) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	c.sendingData.Add(1)
	defer func() {
		c.sendingData.Done()
	}()
	return c.conn.Write(b)
}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}
