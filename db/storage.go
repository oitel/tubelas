package db

import (
	"context"
	"embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/oitel/tubelas/message"
	"github.com/pressly/goose/v3"
)

type impl struct {
	db        *sqlx.DB
	loadStmt  *sqlx.Stmt
	storeStmt *sqlx.Stmt
}

func newStorage() Storage {
	return &impl{}
}

//go:embed migrations/*.sql
var migrations embed.FS

func (s *impl) Open(ctx context.Context, dbstring string) error {
	var err error
	s.db, err = sqlx.Open("postgres", dbstring)
	if err != nil {
		return err
	}
	if err = s.db.PingContext(ctx); err != nil {
		return err
	}

	s.loadStmt, err = s.db.PreparexContext(ctx, `
		SELECT
			id,
			ts,
			text
		FROM messages
		ORDER BY id DESC
		LIMIT $1
	`)
	if err != nil {
		return err
	}

	s.storeStmt, err = s.db.PreparexContext(ctx, `
		INSERT INTO messages(
			ts, text
		) VALUES (
			$1, $2
		) RETURNING id
	`)
	if err != nil {
		return err
	}

	// migrate database
	goose.SetBaseFS(migrations)
	if err := goose.Up(s.db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}

func (s *impl) Close() error {
	return s.db.Close()
}

func (s *impl) Load(ctx context.Context, maxCount uint) ([]message.Message, error) {
	rows, err := s.loadStmt.QueryxContext(ctx, maxCount)
	if err != nil {
		return nil, err
	}

	msgs := []message.Message{}
	for rows.Next() {
		var msg message.Message
		if err := rows.StructScan(&msg); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// reverse list
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nil
}

func (s *impl) Store(ctx context.Context, msg message.Message) (message.Message, error) {
	err := s.storeStmt.GetContext(ctx, &msg.ID, msg.Timestamp, msg.Text)
	return msg, err
}
