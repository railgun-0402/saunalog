package user

import (
	"strings"
	"testing"
	"time"
)

func validParams() User {
	return User{
		ID:         "u-123",
		Name:       "テスト太郎",
		Email:      "test@example.com",
		Password:   "hashed-password",
		Gender:     "male",
		Age:        28,
		Prefecture: "Tokyo",
	}
}

func TestNewUser_Success(t *testing.T) {
	in := validParams()

	before := time.Now()
	u, err := NewUser(in)
	after := time.Now()

	if err != nil {
		t.Fatalf("NewUser() error = %v, want nil", err)
	}
	if u == nil {
		t.Fatalf("NewUser() returned nil user")
	}

	// 値が反映されていること
	if u.ID != in.ID ||
		u.Name != in.Name ||
		u.Email != in.Email ||
		u.Password != in.Password ||
		u.Gender != in.Gender ||
		u.Age != in.Age ||
		u.Prefecture != in.Prefecture {
		t.Errorf("NewUser() fields mismatch\ngot:  %+v\nwant: %+v", *u, in)
	}

	// CreatedAt が現在時刻で設定されていること（before〜after の範囲内）
	if u.CreatedAt.Before(before) || u.CreatedAt.After(after) {
		t.Errorf("CreatedAt out of expected range: got=%v, before=%v, after=%v",
			u.CreatedAt, before, after)
	}
}

func TestNewUser_Validation_MissingFields(t *testing.T) {
	type fields struct {
		name       bool
		gender     bool
		email      bool
		prefecture bool
	}

	cases := []struct {
		name         string
		miss         fields
		wantContains string
	}{
		{
			name:         "missing name",
			miss:         fields{name: true},
			wantContains: "名前 は必須です",
		},
		{
			name:         "missing gender",
			miss:         fields{gender: true},
			wantContains: "性別 は必須です",
		},
		{
			name:         "missing email",
			miss:         fields{email: true},
			wantContains: "メールアドレス は必須です",
		},
		{
			name:         "missing prefecture",
			miss:         fields{prefecture: true},
			wantContains: "都道府県 は必須です",
		},
		{
			name:         "multiple missing (order: 名前・性別・メールアドレス・都道府県)",
			miss:         fields{name: true, gender: true, email: true, prefecture: true},
			wantContains: "名前・性別・メールアドレス・都道府県 は必須です",
		},
		{
			name:         "missing name and email (keeps defined order)",
			miss:         fields{name: true, email: true},
			wantContains: "名前・メールアドレス は必須です",
		},
	}

	for _, tc := range cases {
		/*
		 * 第一引数：サブテストの名前（ここでは tc.name）
		 * 第二引数：そのサブテストを実行する関数
		 */
		t.Run(tc.name, func(t *testing.T) {
			in := validParams()
			if tc.miss.name {
				in.Name = ""
			}
			if tc.miss.gender {
				in.Gender = ""
			}
			if tc.miss.email {
				in.Email = ""
			}
			if tc.miss.prefecture {
				in.Prefecture = ""
			}

			u, err := NewUser(in)
			if err == nil {
				t.Fatalf("NewUser() error = nil, want error")
			}
			if u != nil {
				t.Fatalf("NewUser() user = %#v, want nil", u)
			}
			if !strings.Contains(err.Error(), tc.wantContains) {
				t.Errorf("error message mismatch\n got: %q\nwant: %q", err.Error(), tc.wantContains)
			}
		})
	}
}

func TestNewUser_Validation_AgeNegative(t *testing.T) {
	in := validParams()
	in.Age = -1

	u, err := NewUser(in)
	if err == nil {
		t.Fatalf("NewUser() error = nil, want error")
	}
	if u != nil {
		t.Fatalf("NewUser() user = %#v, want nil", u)
	}
	want := "年齢は0以上で設定してください"
	if err.Error() != want {
		t.Errorf("error message mismatch\n got: %q\nwant: %q", err.Error(), want)
	}
}
