package store

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	databaseURL string
	dbPool      *pgxpool.Pool
}

func New(url string) *Store {
	return &Store{
		databaseURL: url,
	}
}

func (s *Store) Open() error {
	dbpool, err := pgxpool.Connect(context.Background(), s.databaseURL)
	s.dbPool = dbpool
	return err
}

func (s *Store) Close() {
	s.dbPool.Close()
}
