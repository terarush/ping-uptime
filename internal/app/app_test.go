package app

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"ping-uptime/internal/pkg/logger"
)

// repoRoot returns the absolute path to the repository root, computed
// relative to this test file's own location on disk (internal/app/app_test.go).
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine caller information")
	}
	// thisFile = <root>/internal/app/app_test.go
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

// newTestAppWithDocsRoute creates an App (without running the full
// Initialize(), which has a hard dependency on opening a real database) and
// registers only the Scalar API docs route added in this PR, so the new
// routing behavior can be exercised in isolation.
func newTestAppWithDocsRoute(t *testing.T) *App {
	t.Helper()

	cfg := logger.DefaultConfig()
	cfg.OutputPath = filepath.Join(t.TempDir(), "app.log")

	var emptyStaticFS embed.FS
	a, err := NewApp(&cfg, emptyStaticFS)
	if err != nil {
		t.Fatalf("NewApp returned unexpected error: %v", err)
	}

	a.r = a.SetRouter()
	a.registerDocsRoute()

	return a
}

// TestDocsRoute_ServesHTMLPage verifies that requesting any path under
// /api/docs/ (other than openapi.json) returns the Scalar API reference
// HTML page with the expected script tag pointing at the OpenAPI spec.
func TestDocsRoute_ServesHTMLPage(t *testing.T) {
	a := newTestAppWithDocsRoute(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docs/", nil)
	rec := httptest.NewRecorder()

	a.r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected Content-Type to contain text/html, got %q", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `data-url="/api/docs/openapi.json"`) {
		t.Errorf("expected body to reference /api/docs/openapi.json, got: %s", body)
	}
	if !strings.Contains(body, "Ping Uptime API") {
		t.Errorf("expected body to contain page title, got: %s", body)
	}
	if !strings.Contains(body, "@scalar/api-reference") {
		t.Errorf("expected body to reference the scalar api-reference script, got: %s", body)
	}
}

// TestDocsRoute_ServesHTMLPageForArbitrarySubpath verifies the wildcard
// route serves the same HTML page regardless of the specific sub-path
// requested under /api/docs/, as long as it isn't "openapi.json".
func TestDocsRoute_ServesHTMLPageForArbitrarySubpath(t *testing.T) {
	a := newTestAppWithDocsRoute(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docs/some/random/path", nil)
	rec := httptest.NewRecorder()

	a.r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Ping Uptime API") {
		t.Errorf("expected HTML page body, got: %s", rec.Body.String())
	}
}

// TestDocsRoute_ServesOpenAPIJSON verifies that requesting
// /api/docs/openapi.json serves the docs/swagger.json file from disk rather
// than the HTML page.
func TestDocsRoute_ServesOpenAPIJSON(t *testing.T) {
	a := newTestAppWithDocsRoute(t)
	root := repoRoot(t)

	// c.File resolves paths relative to the process's current working
	// directory at request time, so temporarily switch to the repository
	// root (where docs/swagger.json actually lives) to exercise the real
	// file-serving behavior.
	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("failed to chdir to repo root: %v", err)
	}
	defer func() {
		_ = os.Chdir(prevWD)
	}()

	req := httptest.NewRequest(http.MethodGet, "/api/docs/openapi.json", nil)
	rec := httptest.NewRecorder()

	a.r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "json") {
		t.Errorf("expected Content-Type to contain json, got %q", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"swagger"`) {
		limit := len(body)
		if limit > 200 {
			limit = 200
		}
		t.Errorf("expected served file to contain swagger spec content, got: %s", body[:limit])
	}
}

// TestDocsRoute_OpenAPIJSONMissingFileReturnsError verifies that when the
// swagger.json file is not present on disk relative to the working
// directory, the route does not silently fall back to the HTML page (the
// request path is exactly "openapi.json", so it should always attempt to
// serve the file).
func TestDocsRoute_OpenAPIJSONMissingFileReturnsError(t *testing.T) {
	a := newTestAppWithDocsRoute(t)

	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	// An empty temp dir has no docs/swagger.json file.
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}
	defer func() {
		_ = os.Chdir(prevWD)
	}()

	req := httptest.NewRequest(http.MethodGet, "/api/docs/openapi.json", nil)
	rec := httptest.NewRecorder()

	a.r.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK && strings.Contains(rec.Body.String(), "Ping Uptime API") {
		t.Errorf("expected file-serving branch to be used (and fail), not the HTML fallback page")
	}
	if rec.Code < http.StatusBadRequest {
		t.Errorf("expected an error status code when the spec file is missing, got %d", rec.Code)
	}
}