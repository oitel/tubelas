package db

import "github.com/oitel/tubelas/message"

type Storage interface {
	Open(dbstring string) error
	Close() error

	Load(maxCount uint) ([]message.Message, error)
	Store(msg message.Message) (message.Message, error)
}

func NewStorage() Storage {
	return newStorage()
}

var global Storage

func GlobalInstance() Storage {
	if global == nil {
		global = newStorage()
	}
	return global
}
