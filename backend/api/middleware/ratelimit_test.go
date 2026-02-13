package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		if val := r.Header.Get("X-Status"); val != "" {
			if val == "401" {
				status = http.StatusUnauthorized
			}
		}
		w.WriteHeader(status)
	})

	handler := rl.LimitMiddleware(nextHandler)

	req := httptest.NewRequest("POST", "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Request %d expected 200 OK, got %d", i+1, rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request 6 expected 429 Too Many Requests, got %d", rr.Code)
	}

	rl.ips.Delete("192.168.1.1")

	reqFail := httptest.NewRequest("POST", "/api/auth/login", nil)
	reqFail.RemoteAddr = "192.168.1.1:1234"
	reqFail.Header.Set("X-Status", "401")

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		if val, ok := rl.ips.Load("192.168.1.1"); ok {
			state := val.(*IPState)
			state.limiter = rate.NewLimiter(rate.Inf, 0)
		}

		handler.ServeHTTP(rr, reqFail)
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Failure attempt %d expected 401, got %d", i+1, rr.Code)
		}
	}

	val, ok := rl.ips.Load("192.168.1.1")
	if !ok {
		t.Fatal("IP state not found")
	}
	state := val.(*IPState)
	if state.blockUntil.IsZero() {
		t.Error("Specific IP should be blocked after 5 failures, but blockUntil is zero")
	}

	reqcheck := httptest.NewRequest("POST", "/api/auth/login", nil)
	reqcheck.RemoteAddr = "192.168.1.1:1234"
	rrCheck := httptest.NewRecorder()
	handler.ServeHTTP(rrCheck, reqcheck)

	if rrCheck.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429 after 5 failures, got %d", rrCheck.Code)
	}

	var resp map[string]string
	json.NewDecoder(rrCheck.Body).Decode(&resp)
	if msg, ok := resp["error"]; !ok {
		t.Error("Expected error message in JSON response")
	} else {
		t.Logf("Got error message: %s", msg)
	}
}
