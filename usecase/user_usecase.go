package usecase

import (
	"context"
	"github.com/google/uuid"
	domain "saunalog/domain/user"
)

type UserUsecase struct {
	repo domain.Repository
}

func NewUserCreate(r domain.Repository) *UserUsecase {
	return &UserUsecase{repo: r}
}

type CreateUserOutput struct {
	ID string `json:"id"`
}

func (u *UserUsecase) Create(ctx context.Context, input *domain.User) (CreateUserOutput, error) {
	input.ID = uuid.New().String()

	saved, err := u.repo.CreateUser(ctx, input)
	if err != nil {
		return CreateUserOutput{}, err
	}
	return CreateUserOutput{ID: saved.ID}, nil
}
