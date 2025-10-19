package repositories

import (
	"context"
	"errors"
	"fmt"

	"goomba41/gotasker/internal/repository/db"

	"github.com/jinzhu/copier"
)

// Repository — контракт для работы с пользователями.
// Сервисы зависят от этого интерфейса, а не от sqlc.
type Repository interface {
	Create(ctx context.Context, email, password string) (*db.User, error)
	// GetByID(ctx context.Context, id int64) (*db.User, error)
	GetByEmail(ctx context.Context, email string) (*db.User, error)
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

	result, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in DB: %w", err)
	}

	user := db.User{}
	if err := copier.Copy(&user, &result); err != nil {
		return nil, fmt.Errorf("failed to create user in DB: %w", err)
	}

	return &user, nil
}

// GetByEmail возвращает пользователя по email.
func (r *repository) GetByEmail(ctx context.Context, email string) (*db.User, error) {
	result, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found: %w", errors.New("entity not found"))
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user := db.User{}
	if err := copier.Copy(&user, &result); err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}
