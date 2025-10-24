package userRepository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"goomba41/gotasker/internal/repository/db"
	userRepository "goomba41/gotasker/internal/repository/user"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Подключаемся
	var err error
	testDB, err = sql.Open("pgx", "host=localhost user=gotasker password=gotasker dbname=gotasker port=5432 sslmode=disable TimeZone=Europe/Moscow")
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	// Опционально: накатываем миграции, если база отдельная для тестирования
	// runMigrations(testDB)

	// Запускаем тесты
	code := m.Run()

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	querier := db.New(testDB)
	repository := userRepository.New(querier, testDB)

	_, err := repository.Create(context.Background(), "anton.borodawkin@yandex.ru", "password")
	if err != nil {
		t.Errorf("❌ %v", err)
	}

	if err == nil {
		t.Log("✅ User creation: OK")
	}

}
