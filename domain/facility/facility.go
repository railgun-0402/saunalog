package domain

import (
	"errors"
	"time"
)

type SaunaFacilityID string

type SaunaFacility struct {
	ID        SaunaFacilityID
	Name      string
	Address   string
	Price     int
	ImageURL  string
	SaunaInfo SaunaInfo
	CreatedAt time.Time
}

type SaunaInfo struct {
	Temperature  int
	Water        int
	HasMeal      bool
	HasRestArea  bool
	HasSleepRoom bool
}

func NewSaunaFacility(params SaunaFacility) (*SaunaFacility, error) {
	// nameとageは必須入力
	if params.Name == "" || params.Address == "" {
		return nil, errors.New("名前・住所は必須です")
	}

	if params.Price < 0 {
		return nil, errors.New("値段は正の数が必要になります")
	}

	return &SaunaFacility{
		ID:        params.ID,
		Name:      params.Name,
		Address:   params.Address,
		Price:     params.Price,
		ImageURL:  params.ImageURL,
		SaunaInfo: params.SaunaInfo,
		CreatedAt: time.Now(),
	}, nil
}

func NewSaunaInfo(temperature, water int, hasMeal, hasRestArea, hasSleepRoom bool) (*SaunaInfo, error) {
	if temperature < -20 || temperature > 120 {
		return nil, errors.New("いくらなんでもしんでしまいます")
	}

	if water < 0 || water > 25 {
		return nil, errors.New("いくらなんでも水風呂の意味なさ過ぎます")
	}

	return &SaunaInfo{
		Temperature:  temperature,
		Water:        water,
		HasMeal:      hasMeal,
		HasRestArea:  hasRestArea,
		HasSleepRoom: hasSleepRoom,
	}, nil
}
