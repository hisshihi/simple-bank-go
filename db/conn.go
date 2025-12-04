// Package db харнит подключени к базам данных и логику работы с ними.
package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DBType   string // "postgres", "mysql", "sqlite"
	SSLMode  string // для postgres
}

func InitDB(config DBConfig, models ...interface{}) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch config.DBType {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
		dialector = postgres.Open(dsn)
	// Добавьте другие диалекторы по мере необходимости (MySQL, SQLite и т.д.)
	default:
		return nil, fmt.Errorf("неизвестный тип базы данных: %s", config.DBType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Автомиграция моделей
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return nil, fmt.Errorf("ошибка миграции: %w", err)
		}
		log.Println("✅ Автомиграция завершена успешно")
	}

	log.Println("✅ Подключение к БД установлено")
	return db, nil
}
