package main

import (
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

// TODO: сделать необязательные поля и посмотреть, как это сделать в fmt
type dsn struct {
	host     string
	user     string
	password string
	dbname   string
	port     int
	sslmode  string
	timezone string
}

var dsnObject = dsn{host: "localhost", user: "gotasker",
	password: "gotasker",
	dbname:   "gotasker",
	port:     5432,
	sslmode:  "disable",
	timezone: "Europe/Moscow",
}

func connectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dsnObject.host,
		dsnObject.user,
		dsnObject.password,
		dsnObject.dbname,
		dsnObject.port,
		dsnObject.sslmode,
		dsnObject.timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	applyMigrations(dsnObject)

	return db
}

func applyMigrations(dsnObject dsn) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dsnObject.user,
		dsnObject.password,
		dsnObject.host,
		dsnObject.port,
		dsnObject.dbname,
		dsnObject.sslmode,
	)

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
