package server

import (
	"godis/datastruct/dict"
	"godis/interface/database"
	"godis/interface/redis"
	"time"
)

const (
	dataDictSize = 1 << 16
	ttlDictSize  = 1 << 10
)

type DB struct {
	index  int
	data   *dict.ConcurrentDict // 数据字典
	ttlMap *dict.ConcurrentDict // 过期字典

	deleteCallback database.KeyEventCallback
}

func makeDB() *DB {
	db := &DB{
		data:   dict.MakeConcurrent(dataDictSize),
		ttlMap: dict.MakeConcurrent(ttlDictSize),
	}
	return db
}

func (db *DB) Exec(c redis.Connection, cmdLine [][]byte) redis.Reply {

}

// GetEntity 返回ket对应的entity
func (db *DB) GetEntity(key string) (database.DataEntity, bool) {
	raw, ok := db.data.GetWithLock(key)
	if !ok {
		return nil, false
	}
	if db.IsExpired(key) {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

// IsExpired check whether a key has expired
func (db *DB) IsExpired(key string) bool {
	rawExpireTime, ok := db.ttlMap.Get(key)
	if !ok {
		return false
	}
	expireTime, _ := rawExpireTime.(time.Time)
	expired := time.Now().After(expireTime)
	if expired {
		db.Remove(key)
	}
	return expired
}

func (db *DB) Remove(key string) {
	raw, deleted := db.data.RemoveWithLock(key)
	db.ttlMap.Remove(key)

	if cb := db.deleteCallback; cb != nil {
		var entity *database.DataEntity
		if deleted > 0 {
			entity = raw.(*database.DataEntity)
		}
		cb(db.index, key, entity)
	}
}
