package userRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"goomba41/gotasker/internal/dto"
	"goomba41/gotasker/internal/repository/db"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
)

// Repository — контракт для работы с пользователями.
// Сервисы зависят от этого интерфейса, а не от sqlc.
type Repository interface {
	Create(ctx context.Context, email, password string) (*db.User, error)
	GetByID(ctx context.Context, id int64) (*db.User, error)
	GetByEmail(ctx context.Context, email string) (*db.User, error)
	Update(ctx context.Context, id int64, user dto.UserUpdate) (*db.User, error)
	Patch(ctx context.Context, id int64, patch dto.UserPatch) (*db.User, error)
	Delete(ctx context.Context, id int64) (*db.User, error)
}

type DBRunner interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
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

// Create создаёт нового пользователя.
func (r *repository) Create(ctx context.Context, email, password string) (*db.User, error) {
	params := db.CreateUserParams{
		Email:    email,
		Password: password,
	}

	result, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("User creation failed: %w", err)
	}

	user := db.User{}
	if err := copier.Copy(&user, &result); err != nil {
		return nil, fmt.Errorf("User creation failed: %w", err)
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

// GetById возвращает пользователя по id.
func (r *repository) GetByID(ctx context.Context, id int64) (*db.User, error) {
	result, err := r.queries.GetUserById(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found: %w", errors.New("entity not found"))
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	user := db.User{}
	if err := copier.Copy(&user, &result); err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// Update обновляет всего пользователя полностью
func (r *repository) Update(ctx context.Context, id int64, update dto.UserUpdate) (*db.User, error) {
	params := db.UpdateUserParams{}
	if err := copier.Copy(params, &update); err != nil {
		return nil, fmt.Errorf("failed to update user in DB: %w", err)
	}

	result, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update user in DB: %w", errors.New("entity not found"))
		}
		return nil, fmt.Errorf("failed to update user in DB: %w", err)
	}

	user := db.User{}
	if err := copier.Copy(&user, &result); err != nil {
		return nil, fmt.Errorf("failed to create user in DB: %w", err)
	}

	return &user, nil
}

// Patch обновляет только указанные поля пользователя
func (r *repository) Patch(ctx context.Context, id int64, patch dto.UserPatch) (*db.User, error) {
	update := r.sb.Update("users").Where(squirrel.Eq{"id": id})

	if patch.Email != nil {
		update = update.Set("email", *patch.Email)
	}
	if patch.Password != nil {
		update = update.Set("password", *patch.Password)
	}

	query, args, err := update.Suffix("RETURNING id, email, password, created_at").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build patch query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute patch query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found: %w", sql.ErrNoRows)
	}

	var user db.User
	err = rows.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	// Защита от неожиданного множественного результата
	if rows.Next() {
		return nil, errors.New("expected exactly one row, got more")
	}

	return &user, nil
}

// Delete удаляет пользователя.
func (r *repository) Delete(ctx context.Context, id int64) (*db.User, error) {
	result, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("failed to delete user in DB: %w", errors.New("entity not found"))
		}
		return nil, fmt.Errorf("failed to delete user in DB: %w", err)
	}

	user := db.User{}
	if err := copier.Copy(&user, &result); err != nil {
		return nil, fmt.Errorf("failed to delete user in DB: %w", err)
	}

	return &user, nil
}
