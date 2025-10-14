package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	// "os"

	// "github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// db := connectDB()

	connectDB()
}

type dsn struct {
	host, user, password, dbname, sslmode, timezone string
	port                                            int
}

type DSNType string

// Определяем допустимые значения как константы
const (
	DSNTypeMigrate DSNType = "migrate"
	DSNTypeGorm    DSNType = "gorm"
)

var dsnObject = dsn{
	host:     "localhost",
	port:     5432,
	user:     "gotasker",
	password: "gotasker",
	dbname:   "gotasker",
	sslmode:  "disable",
	timezone: "Europe/Moscow",
}

func connectDB() *gorm.DB {
	dsn, err := buildPostgresDSN(DSNTypeGorm)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	applyMigrations()

	return db
}

func applyMigrations() {
	dsn, err := buildPostgresDSN(DSNTypeMigrate)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsDir := "file://" + filepath.Join(dir, "..", "..", "migrations")

	// Создаём экземпляр migrate
	m, err := migrate.New(migrationsDir, dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Применяем все миграции вверх
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("✅ Migrations: no changes")
	} else {
		log.Println("✅ Migrations applied successfully")
	}
}

func buildPostgresDSN(dsnType DSNType) (string, error) {
	switch dsnType {
	case "gorm":
		// TODO: сделать формирование строки через text/template
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			dsnObject.host,
			dsnObject.user,
			dsnObject.password,
			dsnObject.dbname,
			dsnObject.port,
			dsnObject.sslmode,
			dsnObject.timezone,
		), nil
	case "migrate":
		// TODO: сделать формирование строки через text/template
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			dsnObject.user,
			dsnObject.password,
			dsnObject.host,
			dsnObject.port,
			dsnObject.dbname,
			dsnObject.sslmode,
		), nil
	default:
		return "", errors.New("unknown postgres DSN type")
	}
}
