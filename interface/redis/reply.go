package redis

// 所有的返回值需要实现的接口
type Reply interface {
	ToBytes() []byte
}
