package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stwalsh4118/mirageapi/internal/status"
	"github.com/stwalsh4118/mirageapi/internal/store"
)

func TestEnvironmentController_CreateAndGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("db open failed: %v", err)
	}
	ec := &EnvironmentController{DB: db, Railway: nil}

	r := gin.New()
	v1 := r.Group("/api/v1")
	ec.RegisterRoutes(v1)

	payload := map[string]any{
		"name": "env-a",
		"type": "dev",
	}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/environments", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("expected id in response")
	}
	if s, _ := created["status"].(string); s != status.StatusCreating {
		t.Fatalf("expected status %q, got %q", status.StatusCreating, s)
	}

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/api/v1/environments/"+id, nil))
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}

func TestEnvironmentController_ListAndDestroy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("db open failed: %v", err)
	}
	ec := &EnvironmentController{DB: db, Railway: nil}

	r := gin.New()
	v1 := r.Group("/api/v1")
	ec.RegisterRoutes(v1)

	// Create
	payload := map[string]any{
		"name": "env-b",
		"type": "dev",
	}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/environments", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id, _ := created["id"].(string)

	// List
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, httptest.NewRequest(http.MethodGet, "/api/v1/environments", nil))
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wList.Code)
	}

	// Destroy
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, httptest.NewRequest(http.MethodDelete, "/api/v1/environments/"+id, nil))
	if wDel.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", wDel.Code)
	}
}
