package posts

import (
	"bytes"
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

	// create module and inject in-memory service for tests
	m := New()
	m.service = NewInMemoryStore()
	// seed
	_, _ = m.service.CreatePost("first-post", "First Post", "initial content")

	a.RegisterModule(m)
	a.SetupRoutes()
	return a
}

func parseBody(t *testing.T, b *httptest.ResponseRecorder) map[string]interface{} {
	var body map[string]interface{}
	if err := json.Unmarshal(b.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	return body
}

func TestListPostsHandler(t *testing.T) {
	a := setupApp()
	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	w := httptest.NewRecorder()
	a.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := parseBody(t, w)
	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success true")
	}
	data, _ := body["data"].([]interface{})
	if len(data) < 1 {
		t.Fatalf("expected at least one post in response")
	}
}

func TestCreateGetUpdateDeleteHandlers(t *testing.T) {
	a := setupApp()

	// Create
	payload := map[string]string{"slug": "new-post", "title": "New", "content": "body"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
	body := parseBody(t, w)
	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success true on create")
	}
	data := body["data"].(map[string]interface{})
	if data["slug"] != "new-post" {
		t.Fatalf("expected slug new-post, got %v", data["slug"])
	}

	// Get
	req = httptest.NewRequest(http.MethodGet, "/api/posts/new-post", nil)
	w = httptest.NewRecorder()
	a.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	body = parseBody(t, w)
	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success true on get")
	}

	// Update
	upd := map[string]string{"title": "Updated", "content": "updated content"}
	ub, _ := json.Marshal(upd)
	req = httptest.NewRequest(http.MethodPut, "/api/posts/new-post", bytes.NewBuffer(ub))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	a.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	body = parseBody(t, w)
	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success true on update")
	}
	ud := body["data"].(map[string]interface{})
	if ud["title"] != "Updated" {
		t.Fatalf("update did not change title")
	}

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/api/posts/new-post", nil)
	w = httptest.NewRecorder()
	a.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	body = parseBody(t, w)
	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success true on delete")
	}
}
