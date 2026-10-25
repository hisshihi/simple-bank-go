package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConns        = 10
	minConns        = 2
	maxConnLifetime = 10 * time.Minute
	maxConnIdleTime = 5 * time.Minute
)

type DB struct {
	pool *pgxpool.Pool
}

func (d *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return d.pool.Begin(ctx)
}

func (d *DB) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return d.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func New(dsn string) (*DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn is empty")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	config.MaxConns = maxConns
	config.MinConns = minConns
	config.MaxConnLifetime = maxConnLifetime
	config.MaxConnIdleTime = maxConnIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ping pool: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Query выполняет SQL-запрос, который возвращает набор строк.
// Принимает контекст запроса ctx, текст SQL и позиционные параметры args.
// Возвращает pgx.Rows для чтения результата и ошибку выполнения (если произошла).
// Обязательно закрыть тело после выполнения Query запроса.
func (d *DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, sql, args...)
}

// QueryRow выполняет SQL-запрос, который ожидает одну строку результата.
// Принимает контекст запроса ctx, текст SQL и позиционные параметры args.
// Возвращает pgx.Row, из которого можно считать значения; фактическая ошибка
// будет доступна при вызове Scan.
func (d *DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, sql, args...)
}

// Exec выполняет SQL-команду, которая не возвращает строки (например, INSERT/UPDATE/DELETE).
// Принимает контекст запроса ctx, текст SQL и позиционные параметры args.
// Возвращает pgconn.CommandTag (метаданные о выполненной команде) и ошибку выполнения (если произошла).
func (d *DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, sql, args...)
}

func (d *DB) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return d.pool.SendBatch(ctx, b)
}

func (d *DB) Close(_ context.Context) error {
	d.pool.Close()
	return nil
}

func (d *DB) Pool() *pgxpool.Pool {
	return d.pool
}
