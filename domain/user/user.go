package user

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID         string
	Name       string
	Email      string // Unique
	Password   string // Hash
	Gender     string // M/F/Others
	Age        int
	Prefecture string // "Tokyo"など
	CreatedAt  time.Time
}

func NewUser(params User) (*User, error) {
	var missingFields []string

	if params.Name == "" {
		missingFields = append(missingFields, "名前")
	}
	if params.Gender == "" {
		missingFields = append(missingFields, "性別")
	}
	if params.Email == "" {
		missingFields = append(missingFields, "メールアドレス")
	}
	if params.Prefecture == "" {
		missingFields = append(missingFields, "都道府県")
	}
	// バリデーションチェックを複数でまとめ
	if len(missingFields) > 0 {
		return nil, fmt.Errorf("%s は必須です", strings.Join(missingFields, "・"))
	}

	if params.Age < 0 {
		return nil, errors.New("年齢は0以上で設定してください")
	}

	return &User{
		ID:         params.ID,
		Name:       params.Name,
		Email:      params.Email,
		Password:   params.Password,
		Gender:     params.Gender,
		Age:        params.Age,
		Prefecture: params.Prefecture,
		CreatedAt:  time.Now(),
	}, nil
}
