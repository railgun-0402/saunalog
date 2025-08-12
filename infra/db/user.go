package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	domain "saunalog/domain/user"
	"saunalog/usecase/repository"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (name, email, gender, age, password, prefecture, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`

	res, err := r.db.ExecContext(ctx, query,
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

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, gender, age, password, prefecture, created_at
		FROM users WHERE email = ?
	`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
