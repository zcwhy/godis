package dict

import (
	"math"
	"sync"
	"sync/atomic"
)

/*
	分片实现并发hash
*/

type ConcurrentDict struct {
	table      []*shard
	count      int32
	shardCount int
}

type shard struct {
	m     map[string]interface{}
	mutex sync.RWMutex
}

// 计算最接近param的2的幂次数
func computeCapacity(param int) (size int) {
	if param <= 16 {
		return 16
	}
	// 扩散bit位上的1
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	}
	return n + 1
}

func MakeConcurrent(shardCount int) *ConcurrentDict {
	shardCount = computeCapacity(shardCount)
	table := make([]*shard, shardCount)
	for i := 0; i < shardCount; i++ {
		table[i] = &shard{
			m: make(map[string]interface{}),
		}
	}
	d := &ConcurrentDict{
		count:      0,
		table:      table,
		shardCount: shardCount,
	}
	return d
}

// hash质数
const prime32 = uint32(16777619)

// hash算法
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// 根据hashCode定位shard，当n为2的整数幂时 h % n == (n - 1) & h
func (dict *ConcurrentDict) spread(hashCode uint32) uint32 {
	if dict == nil {
		panic("dict is nil")
	}
	tableSize := uint32(len(dict.table))
	return (tableSize - 1) & hashCode
}

func (dict *ConcurrentDict) getShard(index uint32) *shard {
	if dict == nil {
		panic("dict is nil")
	}
	return dict.table[index]
}

func (dict *ConcurrentDict) Get(key string) (val interface{}, exists bool) {
	if dict == nil {
		panic("dict is nil")
	}

	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	s := dict.getShard(index)
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, exists = s.m[key]
	return
}

// GetWithLock 调用方持有锁的时候调用该方法
func (dict *ConcurrentDict) GetWithLock(key string) (val interface{}, exists bool) {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	s := dict.getShard(index)
	val, exists = s.m[key]
	return
}

func (dict *ConcurrentDict) Len() int {
	if dict == nil {
		panic("dict is nil")
	}
	return int(atomic.LoadInt32(&dict.count))
}

func (dict *ConcurrentDict) Put(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	s := dict.getShard(index)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.m[key]; ok {
		s.m[key] = val
		return 0
	}
	dict.addCount()
	s.m[key] = val
	return 1
}

func (dict *ConcurrentDict) Remove(key string) (val interface{}, result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	s := dict.getShard(index)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if val, ok := s.m[key]; ok {
		delete(s.m, key)
		dict.decreaseCount()
		return val, 1
	}
	return nil, 0
}

func (dict *ConcurrentDict) RemoveWithLock(key string) (val interface{}, result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	s := dict.getShard(index)

	if val, ok := s.m[key]; ok {
		delete(s.m, key)
		dict.decreaseCount()
		return val, 1
	}
	return val, 0
}

func (dict *ConcurrentDict) decreaseCount() int32 {
	return atomic.AddInt32(&dict.count, -1)
}

func (dict *ConcurrentDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}
