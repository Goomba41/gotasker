package database

import (
	"errors"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dsn struct {
	Host, User, Password, DBName, SslMode, TimeZone string
	Port                                            int
}

type dsnType string

// Определяем допустимые значения как константы
const (
	dsnTypeMigrate dsnType = "migrate"
	dsnTypeGorm    dsnType = "gorm"
)

var dsnObject = dsn{
	Host:     "localhost",
	Port:     5432,
	User:     "gotasker",
	Password: "gotasker",
	DBName:   "gotasker",
	SslMode:  "disable",
	TimeZone: "Europe/Moscow",
}

func Connect() *gorm.DB {
	dsn, err := buildPostgresDSN(dsnTypeGorm)

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
	dsn, err := buildPostgresDSN(dsnTypeMigrate)
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

func buildPostgresDSN(dsnType dsnType) (string, error) {
	var buf strings.Builder
	dsnTemplate := `host={{.Host}} user={{.User}} password={{.Password}} dbname={{.DBName}} port={{.Port}} {{if .SslMode}}sslmode={{.SslMode}}{{end}} {{if .TimeZone}}TimeZone={{.TimeZone}}{{end}}`

	switch dsnType {
	case "gorm":
		break
	case "migrate":
		dsnTemplate = `postgres://{{.User}}:{{.Password}}@{{.Host}}:{{.Port}}/{{.DBName}}{{if .SslMode}}?sslmode={{.SslMode}}{{end}}`

	default:
		return "", errors.New("unknown DSN type")
	}

	tmpl, err := template.New("dsn").Parse(dsnTemplate)
	if err != nil {
		log.Fatalf("Failed to build database %s dsn: %v", dsnType, err)
	}

	err = tmpl.Execute(&buf, dsnObject)
	if err != nil {
		log.Fatalf("Failed to build database %s dsn: %v", dsnType, err)
	}

	return buf.String(), nil
}
