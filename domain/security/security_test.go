package security_test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"saunalog/domain/security"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	// Set Sample Env
	t.Setenv("JWT_SECRET_KEY", "test-secret")

	userID := uint(123)
	role := "admin"
	kind := security.TokenKind("access")

	// access token sample
	tokenStr, _, err := security.GenerateJWT(userID, role, 15*time.Minute, kind, "api")
	fmt.Println(tokenStr)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &security.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return security.GetSecret(), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*security.Claims)
	if !ok || !token.Valid {
		t.Fatal("failed to parse claims or token invalid")
	}

	if claims.Kind != kind {
		t.Errorf("expected kind %s, got %s", kind, claims.Kind)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("token already expired")
	}
}
