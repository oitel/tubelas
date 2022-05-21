package db

import (
	"context"

	"github.com/oitel/tubelas/message"
)

type Storage interface {
	Open(ctx context.Context, dbstring string) error
	Close() error

	Load(ctx context.Context, maxCount uint) ([]message.Message, error)
	Store(ctx context.Context, msg message.Message) (message.Message, error)

	MaxConnCount() int64
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
