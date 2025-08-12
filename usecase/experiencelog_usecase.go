package usecase

import (
	"context"
	domain "saunalog/domain/experience_log"
)

type ExperienceLogUseCase interface {
	CreateExperienceLog(ctx context.Context, log *domain.ExperienceLog) error
}

// TODO: 実装構造体（今はMockで動かす）
type experienceLogUsecaseImpl struct{}

func NewExperienceLogUseCase() ExperienceLogUseCase {
	return &experienceLogUsecaseImpl{}
}

func (u *experienceLogUsecaseImpl) CreateExperienceLog(ctx context.Context, log *domain.ExperienceLog) error {
	// TODO: DB Impl
	// TODO: Transfer Impl File to Repository
	return nil
}
