package database

import (
	"errors"
	"fmt"
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

	"Goomba41/gotasker/pkg/configuration"
)

type dsnType string

// Определяем допустимые значения как константы
const (
	dsnTypeMigrate dsnType = "migrate"
	dsnTypeGorm    dsnType = "gorm"
)

var dsnObject configuration.DatabaseConfig

func SetConfig(dsn configuration.DatabaseConfig) error {

	if dsn == (configuration.DatabaseConfig{}) {
		return fmt.Errorf("variable is emtpy")
	}

	dsnObject = dsn
	return nil
}

func Connect() (*gorm.DB, error) {
	dsn, err := buildPostgresDSN(dsnTypeGorm)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := applyMigrations(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return db, nil
}

func applyMigrations() error {
	dsn, err := buildPostgresDSN(dsnTypeMigrate)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsDir := "file://" + filepath.Join(dir, "..", "..", "migrations")

	// Создаём экземпляр migrate
	m, err := migrate.New(migrationsDir, dsn)
	if err != nil {
		return fmt.Errorf("Failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Применяем все миграции вверх
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Migration failed: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("✅ Migrations: no changes")
	} else {
		log.Println("✅ Migrations applied successfully")
	}

	return nil
}

func buildPostgresDSN(dsnType dsnType) (string, error) {
	if dsnObject == (configuration.DatabaseConfig{}) {
		return "", fmt.Errorf("%s dsn: %s", dsnType, "configuration is not set because config is empty")
	}

	var buf strings.Builder
	dsnTemplate := `host={{.Host}} user={{.User}} password={{.Password}} dbname={{.DbName}} port={{.Port}} {{if .SslMode}}sslmode={{.SslMode}}{{end}} {{if .TimeZone}}TimeZone={{.TimeZone}}{{end}}`

	switch dsnType {
	case "gorm":
		break
	case "migrate":
		dsnTemplate = `postgres://{{.User}}:{{.Password}}@{{.Host}}:{{.Port}}/{{.DbName}}{{if .SslMode}}?sslmode={{.SslMode}}{{end}}`

	default:
		return "", errors.New("unknown DSN type")
	}

	tmpl, err := template.New("dsn").Parse(dsnTemplate)
	if err != nil {
		return "", fmt.Errorf("%s dsn: %w", dsnType, err)
	}

	err = tmpl.Execute(&buf, dsnObject)
	if err != nil {
		return "", fmt.Errorf("%s dsn: %w", dsnType, err)
	}

	return buf.String(), nil
}
