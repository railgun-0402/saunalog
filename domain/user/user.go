package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type UserID string

type User struct {
	ID         UserID
	Name       string
	Email      string // Unique
	Password   string // Hash
	Gender     string // M/F/Others
	Age        int
	Prefecture string // "Tokyo"など
	CreatedAt  time.Time
}

func NewUser(id UserID, name, gender, password, email string, age int, prefecture string) (*User, error) {
	var missingFields []string

	if name == "" {
		missingFields = append(missingFields, "名前")
	}
	if gender == "" {
		missingFields = append(missingFields, "性別")
	}
	if email == "" {
		missingFields = append(missingFields, "メールアドレス")
	}
	if prefecture == "" {
		missingFields = append(missingFields, "都道府県")
	}
	// バリデーションチェックを複数でまとめ
	if len(missingFields) > 0 {
		return nil, fmt.Errorf("%s は必須です", strings.Join(missingFields, "・"))
	}

	if age < 0 {
		return nil, errors.New("年齢は0以上で設定してください")
	}

	return &User{
		ID:         id,
		Name:       name,
		Email:      email,
		Password:   password,
		Gender:     gender,
		Age:        age,
		Prefecture: prefecture,
		CreatedAt:  time.Now(),
	}, nil
}
