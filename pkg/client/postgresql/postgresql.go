package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/alexeybs90/go_bus_routes/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Storage struct {
	client Client
}

func (s *Storage) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return s.client.Exec(ctx, sql, arguments...)
}

func (s *Storage) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return s.client.Query(ctx, sql, args...)
}

func (s *Storage) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return s.client.QueryRow(ctx, sql, args...)
}

func (s *Storage) Begin(ctx context.Context) (pgx.Tx, error) {
	return s.client.Begin(ctx)
}

func NewClient(ctx context.Context, cfg config.Storage) (*Storage, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return &Storage{client: dbpool}, nil
}
