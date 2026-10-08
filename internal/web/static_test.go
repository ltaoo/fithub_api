package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStaticAppsAndFallbacks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	static_dir := t.TempDir()
	for _, app_name := range []string{"admin", "h5"} {
		if err := os.MkdirAll(filepath.Join(static_dir, app_name, "assets"), 0755); err != nil {
			t.Fatal(err)
		}
		for filename, content := range map[string]string{"index.html": "<html>" + app_name + "</html>", "assets/app.js": app_name + "-script"} {
			if err := os.WriteFile(filepath.Join(static_dir, app_name, filename), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	router := gin.New()
	router.GET("/api/ping", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "ok"}) })
	Register(router, static_dir)
	cases := []struct {
		method, url, accept string
		status              int
		body                string
	}{
		{"GET", "/", "", 200, "h5"},
		{"GET", "/workout_plan_profile?id=1", "text/html", 200, "h5"},
		{"GET", "/admin/", "", 200, "admin"},
		{"GET", "/admin/login", "text/html", 200, "admin"},
		{"HEAD", "/admin/login", "text/html", 200, ""},
		{"GET", "/assets/app.js", "", 200, "h5-script"},
		{"GET", "/admin/assets/app.js", "", 200, "admin-script"},
		{"GET", "/assets/missing.js", "", 404, ""},
		{"GET", "/admin/assets/missing", "", 404, ""},
		{"GET", "/api/missing", "text/html", 404, ""},
		{"POST", "/admin/login", "text/html", 404, ""},
		{"GET", "/not-a-page", "application/json", 404, ""},
		{"GET", "/.env", "", 404, ""},
		{"GET", "/admin/../.env", "", 404, ""},
		{"GET", "/admin/assets/", "", 404, ""},
		{"GET", "/api/ping", "", 200, "ok"},
	}
	for _, test_case := range cases {
		t.Run(test_case.method+test_case.url, func(t *testing.T) {
			request := httptest.NewRequest(test_case.method, test_case.url, nil)
			request.Header.Set("Accept", test_case.accept)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code != test_case.status || !strings.Contains(recorder.Body.String(), test_case.body) {
				t.Fatalf("status=%d body=%q", recorder.Code, recorder.Body.String())
			}
			if test_case.method == "HEAD" && recorder.Body.Len() != 0 {
				t.Fatal("HEAD returned a body")
			}
		})
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/admin?redirect=1", nil))
	if recorder.Code != http.StatusPermanentRedirect || recorder.Header().Get("Location") != "/admin/?redirect=1" {
		t.Fatal("admin redirect lost path or query")
	}
	// A symlink outside the application tree must not expose server-side files.
	outside_file := filepath.Join(t.TempDir(), "secret.txt")
	os.WriteFile(outside_file, []byte("secret"), 0644)
	if err := os.Symlink(outside_file, filepath.Join(static_dir, "h5", "leak.txt")); err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/leak.txt", nil))
	if recorder.Code != 404 {
		t.Fatal("served external symlink")
	}
}

func TestUnbuiltFrontendDoesNotInterceptAPI(t *testing.T) {
	router := gin.New()
	Register(router, t.TempDir())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))
	if recorder.Code != 503 {
		t.Fatal(recorder.Code)
	}
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/api/missing", nil))
	if recorder.Code != 404 {
		t.Fatal(recorder.Code)
	}
}
