package security

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"sync"
	"time"
)

// TokenKind Access or Refresh
type TokenKind string

const (
	Issuer           = "saunalog"
	ClockSkew        = 30 * time.Second
	TokenKindAccess  = TokenKind("access")
	TokenKindRefresh = TokenKind("refresh")
)

var (
	secret     []byte
	secretOnce sync.Once
)

func GetSecret() []byte {
	secretOnce.Do(func() {
		secret = []byte(os.Getenv("JWT_SECRET_KEY"))
	})
	return secret
}

type Claims struct {
	Kind TokenKind `json:"kind"`
	Role string    `json:"role, omitempty`
	jwt.RegisteredClaims
}

func newJTI() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func GenerateJWT(userID uint, role string, ttl time.Duration, kind TokenKind, audience ...string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)
	jti, err := newJTI()
	if err != nil {
		return "", time.Time{}, err
	}

	sub := strconv.FormatUint(uint64(userID), 10)

	// Create payload
	claims := &Claims{
		Kind: kind,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   sub,
			Audience:  audience,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-ClockSkew)),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        jti,
		},
	}

	// Create token And sign
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signature
	tokenString, err := token.SignedString(GetSecret())
	return tokenString, exp, err
}

// TODO: Error Handle

// How To Use
// Access 15Min, Refresh 14Days
//access,  accessExp,  err := security.GenerateJWT(uid, role, 15*time.Minute, security.TokenKindAccess, "api")
//refresh, refreshExp, err := security.GenerateJWT(uid, "",   14*24*time.Hour, security.TokenKindRefresh)
