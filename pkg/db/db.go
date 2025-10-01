package db

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

type DBOptions struct {
	TimeoutSec      int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // maximum connection lifetime in seconds
	ConnMaxIdleTime int // maximum connection idle time in seconds
	EnableTracing   bool
}

type DbOrTx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewDbWithWait(connectionString string, dbOptions DBOptions) (*pgxpool.Pool, error) {
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	logger.GlobalSugared().Infof(`Connecting to database, timeout after %d seconds`, dbOptions.TimeoutSec)
	timeout := time.Duration(dbOptions.TimeoutSec) * time.Second

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	err = db.Ping(ctx)
	for err != nil {
		nextErr := db.Ping(ctx)
		if errors.Is(nextErr, context.DeadlineExceeded) || errors.Is(nextErr, context.Canceled) {
			return nil, err
		}
		err = nextErr
	}

	return db, nil
}
