package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type stubJWT struct {
	validateResult bool
	validateErr    error
	claims         map[string]interface{}
	parseErr       error
}

func (s *stubJWT) GenerateToken(data map[string]interface{}) (string, error) {
	return "", nil
}

func (s *stubJWT) ValidateToken(token string) (bool, error) {
	return s.validateResult, s.validateErr
}

func (s *stubJWT) ParseToken(token string) (map[string]interface{}, error) {
	return s.claims, s.parseErr
}

func TestAuthRejectsMissingAuthorizationHeader(t *testing.T) {
	InitializeAuth(&stubJWT{})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	err := Auth(func(c echo.Context) error {
		called = true
		return nil
	})(c)
	if err != nil {
		t.Fatalf("Auth returned an unexpected error: %v", err)
	}

	if called {
		t.Fatal("expected next handler not to be called")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthRejectsInvalidToken(t *testing.T) {
	InitializeAuth(&stubJWT{
		validateErr: errors.New("invalid token"),
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Auth(func(c echo.Context) error { return nil })(c)
	if err != nil {
		t.Fatalf("Auth returned an unexpected error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthSetsUserClaimsAndCallsNext(t *testing.T) {
	expectedClaims := map[string]interface{}{
		"user_id": float64(7),
		"role":    "admin",
	}
	InitializeAuth(&stubJWT{
		validateResult: true,
		claims:         expectedClaims,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	err := Auth(func(c echo.Context) error {
		called = true
		if got := c.Get("user"); got == nil {
			t.Fatal("expected user claims to be attached to context")
		}
		return c.NoContent(http.StatusNoContent)
	})(c)
	if err != nil {
		t.Fatalf("Auth returned an unexpected error: %v", err)
	}

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

func TestAdminRejectsNonAdminUsers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", map[string]interface{}{"role": "user"})

	err := Admin(func(c echo.Context) error { return nil })(c)
	if err != nil {
		t.Fatalf("Admin returned an unexpected error: %v", err)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
}

func TestAdminAllowsAdminUsers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", map[string]interface{}{"role": "admin"})

	called := false
	err := Admin(func(c echo.Context) error {
		called = true
		return c.NoContent(http.StatusNoContent)
	})(c)
	if err != nil {
		t.Fatalf("Admin returned an unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected next handler to be called")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
