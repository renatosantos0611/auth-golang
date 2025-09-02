package repositories

import (
	"auth-golang/internal/database"
	"auth-golang/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UsersRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *database.Service) *UserRepository {
	return &UserRepository{
		db: db.DB,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, name, username, email, password, refresh_token, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.RefreshToken,
		int(user.Role),
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	query := `
		SELECT id, name, username, email, password, refresh_token, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	var roleInt int

	err = r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RefreshToken,
		&roleInt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.Role = models.Role(roleInt)
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, username, email, password, refresh_token, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	var roleInt int

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RefreshToken,
		&roleInt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.Role = models.Role(roleInt)
	return &user, nil
}

func (r *UserRepository) FindByRefreshToken(ctx context.Context, token string) (*models.User, error) {
	query := `
		SELECT id, name, username, email, password, refresh_token, role, created_at, updated_at
		FROM users
		WHERE refresh_token = $1
	`

	var user models.User
	var roleInt int

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RefreshToken,
		&roleInt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.Role = models.Role(roleInt)
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if user.ID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	user.UpdatedAt = time.Now()

	query := `
		UPDATE users 
		SET name = $2, username = $3, email = $4, password = $5, refresh_token = $6, role = $7, updated_at = $8
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.RefreshToken,
		int(user.Role),
		user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, nil
	}

	return user, nil
}
