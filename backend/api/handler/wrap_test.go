package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"setlist/api/apierror"
	"testing"
)

// --- Wrap ---

func TestWrap_HandlerReturnsNil(t *testing.T) {
	h := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	Wrap(h)(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestWrap_HandlerReturnsAppError_404(t *testing.T) {
	h := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return apierror.NotFound("Resource")
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	Wrap(h)(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}
	if body["code"] != apierror.ErrNotFound {
		t.Errorf("expected code %q, got %q", apierror.ErrNotFound, body["code"])
	}
}

func TestWrap_HandlerReturnsAppError_400(t *testing.T) {
	h := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return apierror.InvalidRequest("champ manquant")
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	Wrap(h)(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWrap_HandlerReturnsGenericError_500(t *testing.T) {
	h := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("unexpected failure")
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	Wrap(h)(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}
	if body["code"] != apierror.ErrInternal {
		t.Errorf("expected code %q, got %q", apierror.ErrInternal, body["code"])
	}
}
