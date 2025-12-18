package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"softeng-platform/internal/model"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
}

type userRepository struct {
	db *Database
}

func NewUserRepository(db *Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, nickname, email, password, avatar, description, face_photo, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Nickname,
		user.Email,
		user.Password,
		user.Avatar,
		user.Description,
		user.FacePhoto,
		user.Role,
		time.Now(),
		time.Now(),
	)
	if err != nil {

		return fmt.Errorf("failed to create user: %v", err)
	}

	// 获取自增ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}
	user.ID = int(id)


	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `

		SELECT id, username, nickname, email, password, avatar, description, face_photo, role, created_at, updated_at

		FROM users WHERE id = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Description,
		&user.FacePhoto,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, nickname, email, password, avatar, description, face_photo, role, created_at, updated_at
		FROM users WHERE username = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Description,
		&user.FacePhoto,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %v", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, nickname, email, password, avatar, description, face_photo, role, created_at, updated_at
		FROM users WHERE email = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Description,
		&user.FacePhoto,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET nickname = ?, avatar = ?, description = ?, face_photo = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Nickname,
		user.Avatar,
		user.Description,
		user.FacePhoto,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`


	result, err := r.db.ExecContext(ctx, query, hashedPassword, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil

}
