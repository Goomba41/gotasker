package taskRepository

import (
	"context"
	"database/sql"

	// "errors"
	// "fmt"

	"goomba41/gotasker/internal/repository/db"
	// "goomba41/gotasker/internal/dto"

	"github.com/Masterminds/squirrel"
	// "github.com/jinzhu/copier"
)

// Repository — контракт для работы с пользователями.
// Сервисы зависят от этого интерфейса, а не от sqlc.
type Repository interface {
	// Create(ctx context.Context, email, password string) (*db.User, error)
	// GetByID(ctx context.Context, id int64) (*db.User, error)
	// GetByEmail(ctx context.Context, email string) (*db.User, error)
	// Update(ctx context.Context, id int64, user dto.UserUpdate) (*db.User, error)
	// Patch(ctx context.Context, id int64, patch dto.UserPatch) (*db.User, error)
	// Delete(ctx context.Context, id int64) (*db.User, error)
}

type DBRunner interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// repository — реализация Repository с использованием sqlc.
type repository struct {
	queries *db.Queries
	db      DBRunner
	sb      squirrel.StatementBuilderType
}

// NewUserRepository создаёт новый репозиторий для работы с пользователями.
func New(queries *db.Queries, db DBRunner) Repository {
	return &repository{
		queries: queries,
		db:      db,
		sb:      squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
