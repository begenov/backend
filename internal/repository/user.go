package repository

import (
	"context"
	"database/sql"

	"github.com/begenov/backend/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, arg domain.CreateUserParams) (domain.User, error) {
	stmt := `INSERT INTO users (username, hashed_password, full_name, email) 
	VALUES ($1, $2, $3, $4) 
	RETURNING username, hashed_password, full_name, email, password_changed_at, created_at`
	row := r.db.QueryRowContext(ctx, stmt, arg.Username, arg.HashedPassword, arg.FullName, arg.Email)
	var i domain.User
	if err := row.Scan(&i.Username, &i.HashedPassword, &i.FullName, &i.Email, &i.PasswordChangedAt, &i.CreatedAt); err != nil {
		return domain.User{}, err
	}
	return i, nil
}

func (r *UserRepo) GetUser(ctx context.Context, username string) (domain.User, error) {
	stmt := `SELECT username, hashed_password, full_name, email, password_changed_at, created_at FROM users
	WHERE username = $1 LIMIT 1`
	row := r.db.QueryRowContext(ctx, stmt, username)
	var i domain.User
	if err := row.Scan(&i.Username, &i.HashedPassword, &i.FullName, &i.Email, &i.PasswordChangedAt, &i.CreatedAt); err != nil {
		return domain.User{}, err
	}
	return i, nil
}
