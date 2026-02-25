package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"setlist/api/middleware"
	"strings"
	"testing"
)

// --- DecodeJSON ---

func TestDecodeJSON_ValidJSON(t *testing.T) {
	type payload struct{ Name string }
	body := `{"name":"test"}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	v, err := DecodeJSON[payload](r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if v.Name != "test" {
		t.Errorf("expected Name 'test', got '%s'", v.Name)
	}
}

func TestDecodeJSON_MalformedJSON(t *testing.T) {
	type payload struct{ Name string }
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	_, err := DecodeJSON[payload](r)
	if err == nil {
		t.Fatal("expected an error for malformed JSON, got nil")
	}
}

func TestDecodeJSON_EmptyBody(t *testing.T) {
	type payload struct{ Name string }
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))
	_, err := DecodeJSON[payload](r)
	if err == nil {
		t.Fatal("expected an error for empty body, got nil")
	}
}

// --- GetIntParam ---

func TestGetIntParam_ValidInt(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/song/42", nil)
	// Simulate r.PathValue by using a real mux
	mux := http.NewServeMux()
	var captured int
	mux.HandleFunc("GET /song/{id}", func(w http.ResponseWriter, req *http.Request) {
		id, err := GetIntParam(req, "id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		captured = id
	})
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if captured != 42 {
		t.Errorf("expected 42, got %d", captured)
	}
}

func TestGetIntParam_NonNumeric(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/song/abc", nil)
	mux := http.NewServeMux()
	var gotError bool
	mux.HandleFunc("GET /song/{id}", func(w http.ResponseWriter, req *http.Request) {
		_, err := GetIntParam(req, "id")
		if err != nil {
			gotError = true
		}
	})
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if !gotError {
		t.Error("expected error for non-numeric param, got nil")
	}
}

// --- GetBandID ---

func TestGetBandID_Present(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(r.Context(), middleware.BandIDKey, 7)
	r = r.WithContext(ctx)

	id, err := GetBandID(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 7 {
		t.Errorf("expected 7, got %d", id)
	}
}

func TestGetBandID_Missing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetBandID(r)
	if err == nil {
		t.Fatal("expected error when BandIDKey is absent, got nil")
	}
}

// --- GetUserID ---

func TestGetUserID_Present(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, 3)
	r = r.WithContext(ctx)

	id, err := GetUserID(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 3 {
		t.Errorf("expected 3, got %d", id)
	}
}

func TestGetUserID_Missing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetUserID(r)
	if err == nil {
		t.Fatal("expected error when UserIDKey is absent, got nil")
	}
}
