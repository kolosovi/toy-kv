package walmanager

import toykv "github.com/kolosovi/toy-kv"

type Insert struct {
	toykv.Record
}

type Delete struct {
	toykv.K
}

type Log struct {
	typ LogType
	Insert
	Delete
}

type LogType uint

const (
	LogTypeInvalid LogType = iota
	LogTypeInsert
	LogTypeDelete
)