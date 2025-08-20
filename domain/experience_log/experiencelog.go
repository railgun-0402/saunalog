package domain

import (
	"errors"
	domain "saunalog/domain/facility"
	"time"
)

type ExperienceID string

type ExperienceLog struct {
	ID              ExperienceID
	UserID          string
	SaunaFacilityID domain.SaunaFacilityID
	Date            time.Time
	CongestionLevel int // 1〜5(混雑)
	CostPerformance int // 1〜5(コスパ)
	TotonoiLevel    int // 読みずらいが整いレベルを1〜5に
	Comment         string
	CreatedAt       time.Time
}

func NewExperienceLog(params ExperienceLog) (*ExperienceLog, error) {
	// 整い度と混雑度は5段階で評価する
	if !isRatingValid(params.TotonoiLevel) {
		return nil, errors.New("整い度は1〜5で指定してください")
	}
	if !isRatingValid(params.CongestionLevel) {
		return nil, errors.New("混雑度は1〜5で指定してください")
	}
	if !isRatingValid(params.CostPerformance) {
		return nil, errors.New("コスパは1〜5で指定してください")
	}

	// 他にもルールがあればここで追加
	return &ExperienceLog{
		ID:              params.ID,
		UserID:          params.UserID,
		SaunaFacilityID: params.SaunaFacilityID,
		Date:            params.Date,
		CongestionLevel: params.CongestionLevel,
		CostPerformance: params.CostPerformance,
		TotonoiLevel:    params.TotonoiLevel,
		Comment:         params.Comment,
		CreatedAt:       time.Now(),
	}, nil
}

func isRatingValid(v int) bool {
	return v >= 1 && v <= 5
}
