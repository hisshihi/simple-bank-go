package repo

import (
	"errors"
	"testing"

	"github.com/hisshihi/simple-bank/core/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB создает и возвращает временную базу данных для тестирования.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Миграция
	if err := database.AutoMigrate(&db.Account{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Cleanup
	t.Cleanup(func() {
		sqlDB, _ := database.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	})

	return database
}

// assertNoError проверяет отсутствие ошибки
func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// assertErrorIs проверяет конкретную ошибку через errors.Is
func assertErrorIs(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("expected error %v, got %v", want, got)
	}
}

func assertError(t testing.TB, err error, expected error) {
	t.Helper()
	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}
