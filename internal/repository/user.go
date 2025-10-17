package repositories

import (
	"context"
	"fmt"

	"goomba41/gotasker/internal/repository/db"
)

// Repository — контракт для работы с пользователями.
// Сервисы зависят от этого интерфейса, а не от sqlc.
type Repository interface {
	Create(ctx context.Context, email, password string) (*db.User, error)
	// GetByID(ctx context.Context, id int64) (*db.User, error)
	// GetByEmail(ctx context.Context, email string) (*db.User, error)
	// Delete(ctx context.Context, id int64) (*db.User, error)
}

// repository — реализация Repository с использованием sqlc.
type repository struct {
	queries *db.Queries
}

// NewUserRepository создаёт новый репозиторий для работы с пользователями.
func NewUserRepository(queries *db.Queries) Repository {
	return &repository{queries: queries}
}

// Create создаёт нового пользователя.
func (r *repository) Create(ctx context.Context, email, password string) (*db.User, error) {
	params := db.CreateUserParams{
		Email:    email,
		Password: password,
	}

	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in DB: %w", err)
	}

	return user, nil
}
