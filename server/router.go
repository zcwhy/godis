package server

import (
	"godis/interface/redis"
	"strings"
)

var cmdTable = make(map[string]*command)

// ExecFunc 命令行的执行函数
type ExecFunc func(db *DB, args [][]byte) redis.Reply

type PreFunc func(args [][]byte) ([]string, []string)

// CmdLine [][]byte,表示一条指令
type CmdLine = [][]byte

type UndoFunc func(db *DB, args [][]byte) []CmdLine

type command struct {
	name string

	prepare  PreFunc  // 命令前置处理函数
	executor ExecFunc // 命令执行函数
	undo     UndoFunc // 在命令真正执行之前会产生一条undo日志，用于rolled back

	// 命令行允许的参数个数，arity < 0表示参数len(args) >= -arity
	// example: arity of 'get' = 2, 'mget' = -2
	arity int
	flags int
}

func registerCommand(name string, executor ExecFunc, prepare PreFunc, rollback UndoFunc, arity int, flags int) *command {
	name = strings.ToLower(name)
	cmd := &command{
		name:     name,
		executor: executor,
		prepare:  prepare,
		undo:     rollback,
		arity:    arity,
		flags:    flags,
	}
	cmdTable[name] = cmd
	return cmd
}

