package redis

// Connection 需要实现的方法
type Connection interface {
	Write([]byte) (int, error)

	GetDBIndex() int
}
