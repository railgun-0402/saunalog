package domain

import (
	"errors"
	"time"
)

type UserID string

type User struct {
	ID         UserID
	Name       string
	Gender     string // M/F/Others
	Age        int
	Prefecture string // "Tokyo"など
	CreatedAt  time.Time
}

func NewUser(id UserID, name, gender string, age int, prefecture string) (*User, error) {
	// nameとageは必須入力
	if name == "" || gender == "" {
		return nil, errors.New("名前・性別は必須です")
	}

	if age < 0 {
		return nil, errors.New("年齢は0以上で設定してください")
	}

	return &User{
		ID:         id,
		Name:       name,
		Gender:     gender,
		Age:        age,
		Prefecture: prefecture,
		CreatedAt:  time.Now(),
	}, nil
}
