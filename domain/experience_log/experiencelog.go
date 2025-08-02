package domain

import (
	"errors"
	domain "saunalog/domain/facility"
	user "saunalog/domain/user"
	"time"
)

type ExperienceID string

type ExperienceLog struct {
	ID              ExperienceID
	UserID          user.UserID
	SaunaFacilityID domain.SaunaFacilityID
	Date            time.Time
	CongestionLevel int // 1〜5(混雑)
	CostPerformance int
	TotonoiLevel    int // 読みずらいが整いレベルを1〜5に
	Comment         string
	CreatedAt       time.Time
}

func NewExperienceLog(userID user.UserID, facilityID domain.SaunaFacilityID, congestion, totonoi, cost int, comment string) (*ExperienceLog, error) {
	// 整い度と混雑度は5段階で評価する
	if totonoi < 1 || totonoi > 5 {
		return nil, errors.New("整い度は1〜5で指定してください")
	}
	if congestion < 1 || congestion > 5 {
		return nil, errors.New("混雑度は1〜5で指定してください")
	}

	// 費用
	if cost < 0 {
		return nil, errors.New("費用は正の整数で設定してください")
	}

	// 他にもルールがあればここで追加
	return &ExperienceLog{
		UserID:          userID,
		SaunaFacilityID: facilityID,
		Date:            time.Now(),
		CongestionLevel: congestion,
		CostPerformance: cost,
		TotonoiLevel:    totonoi,
		Comment:         comment,
		CreatedAt:       time.Now(),
	}, nil
}
