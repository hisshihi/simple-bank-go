package helper

import (
	"context"
	"log"
	"testing"

	"github.com/hisshihi/simple-bank/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testPool *pgxpool.Pool

func TestSetupPool(m *testing.M) (code int) {
	ctx := context.Background()

	container, err := pgcontainer.Run(ctx,
		"postgres:18.3-alpine3.23",
		pgcontainer.WithDatabase("postgres"),
		pgcontainer.WithUsername("postgres"),
		pgcontainer.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2)))
	if err != nil {
		log.Fatalf("не удалось поднять тестовый контейнер: %v", err)
	}
	defer func() {
		if err = container.Terminate(ctx); err != nil {
			log.Fatalf("ошибка при остановке контейнера: %v", err)
		}
	}()

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("не удалось получить подключение: %v", err)
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("не удалось создать пул: %v", err)
	}
	defer testPool.Close()

	return m.Run()
}

func runMigration(ctx context.Context, pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatalf("err to close connection: %v", err)
		}
	}()

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.UpContext(ctx, sqlDB, ".")
}

// NewTx создаёт транзакцию для одного теста.
// После теста делает Rollback — следующий тест получает чистую БД.
// Это быстрее чем TRUNCATE и не требует порядка выполнения тестов.
func NewTx(t *testing.T) pgx.Tx {
	t.Helper()

	tx, err := testPool.Begin(t.Context())
	if err != nil {
		t.Fatalf("не удалось начать транзакцию: %v", err)
	}

	t.Cleanup(func() {
		_ = tx.Rollback(context.Background())
	})

	return tx
}
