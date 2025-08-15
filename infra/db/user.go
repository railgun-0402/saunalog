package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	domain "saunalog/domain/user"
	"saunalog/usecase/repository"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (name, email, gender, age, password, prefecture, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`

	res, err := r.DB.ExecContext(ctx, query,
		user.Name, user.Email, user.Gender, user.Age, user.Password, user.Prefecture,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	user.ID = domain.UserID(id)

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, gender, age, password, prefecture, created_at
		FROM users WHERE email = ?
	`
	u := domain.User{}
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Gender, &u.Age, &u.Password, &u.Prefecture, &u.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
