package server

import (
	"godis/datastruct/list"
	"godis/interface/database"
	"godis/interface/redis"
	"godis/redis/parser"
)

func init() {
	registerCommand("LPush", execLPush, writeFirstKey, undoLPush, -3, flagWrite)

}

// Insert all the specified values at the head of the list stored at key
// return: the len of list after push
func execLPush(db *DB, args [][]byte) redis.Reply {
	key := string(args[0])
	values := args[1:]

	list, _, errReply := db.getOrInitList(key)
	if errReply != nil {
		return errReply
	}

	for _, value := range values {
		list.Insert(0, value)
	}
	return parser.MakeIntReply(int64(list.Len()))
}

func (db *DB) getAsList(key string) (list.List, parser.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	l, ok := entity.(list.List)
	if !ok {
		return nil, &parser.WrongTypeErrReply{}
	}
	return l, nil
}

func (db *DB) getOrInitList(key string) (list.List, bool, parser.ErrorReply) {
	l, errReply := db.getAsList(key)
	if errReply != nil {
		return nil, false, errReply
	}

	isNew := false
	if l == nil {
		l = list.NewQuickList()
		db.PutEntity(key, database.DataEntity(l))
		isNew = true
	}
	return l, isNew, nil
}
