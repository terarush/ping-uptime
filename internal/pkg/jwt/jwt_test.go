package jwt

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateTokenAddsExpirationWhenMissing(t *testing.T) {
	service := NewJWTImpl("secret", 1)

	tokenString, err := service.GenerateToken(map[string]interface{}{
		"user_id": float64(42),
		"role":    "admin",
	})
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	claims, err := service.ParseToken(tokenString)
	if err != nil {
		t.Fatalf("ParseToken returned an error: %v", err)
	}

	if claims["user_id"] != float64(42) {
		t.Fatalf("expected user_id 42, got %v", claims["user_id"])
	}
	if claims["role"] != "admin" {
		t.Fatalf("expected role admin, got %v", claims["role"])
	}
	if _, ok := claims["exp"]; !ok {
		t.Fatal("expected generated token to contain exp claim")
	}
}

func TestGenerateTokenPreservesExplicitExpiration(t *testing.T) {
	service := NewJWTImpl("secret", 99)
	explicitExpiration := time.Now().Add(2 * time.Hour).Unix()

	tokenString, err := service.GenerateToken(map[string]interface{}{
		"exp": explicitExpiration,
	})
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	claims, err := service.ParseToken(tokenString)
	if err != nil {
		t.Fatalf("ParseToken returned an error: %v", err)
	}

	if claims["exp"] != float64(explicitExpiration) {
		t.Fatalf("expected exp %d, got %v", explicitExpiration, claims["exp"])
	}
}

func TestValidateTokenRejectsWrongSignature(t *testing.T) {
	validService := NewJWTImpl("secret", 1)
	otherService := NewJWTImpl("other-secret", 1)

	tokenString, err := validService.GenerateToken(map[string]interface{}{
		"user_id": float64(42),
	})
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	valid, err := otherService.ValidateToken(tokenString)
	if err == nil {
		t.Fatal("expected ValidateToken to return an error for a token signed with a different key")
	}
	if valid {
		t.Fatal("expected ValidateToken to reject a token signed with a different key")
	}
}

func TestParseTokenRejectsMalformedToken(t *testing.T) {
	service := NewJWTImpl("secret", 1)

	_, err := service.ParseToken("not-a-token")
	if err == nil {
		t.Fatal("expected ParseToken to return an error for a malformed token")
	}
	if !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token-related error, got %v", err)
	}
}
