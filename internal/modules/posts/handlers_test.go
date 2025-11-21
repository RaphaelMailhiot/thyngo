package posts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "thyngo/internal/app"

	"github.com/gin-gonic/gin"
)

func setupApp() *app.App {
	gin.SetMode(gin.TestMode)
	a := app.NewApp()
	a.RegisterModule(New())
	a.SetupRoutes()
	return a
}

func TestListPostsHandler(t *testing.T) {
	a := setupApp()
	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	w := httptest.NewRecorder()
	a.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success true")
	}
}
