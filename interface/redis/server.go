package redis

type Server interface {
	Exec(client Connection, cmdLine [][]byte) Reply
	AfterClientClose(c Connection)
	Close()
}
