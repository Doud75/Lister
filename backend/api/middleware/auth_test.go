package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"setlist/auth"

	"github.com/golang-jwt/jwt/v5"
)

const testJWTSecret = "testsecret"

// serveWithAuth runs the JWTAuthUserOnly middleware against a request carrying
// the given Authorization header and reports whether the wrapped handler was
// reached and the resulting HTTP status code.
func serveWithAuth(t *testing.T, authHeader string) (reached bool, status int) {
	t.Helper()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	})

	handler := JWTAuthUserOnly(testJWTSecret)(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	return reached, rec.Code
}

func TestValidateJWT_RejectsNoneAlgorithm(t *testing.T) {
	claims := auth.JWTClaims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed to forge none-alg token: %v", err)
	}

	reached, status := serveWithAuth(t, "Bearer "+token)

	if reached {
		t.Error("handler should not be reached with a none-alg token")
	}
	if status != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, status)
	}
}

func TestValidateJWT_AcceptsValidHS256(t *testing.T) {
	token, err := auth.GenerateJWT(testJWTSecret, 1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	reached, status := serveWithAuth(t, "Bearer "+token)

	if !reached {
		t.Error("handler should be reached with a valid HS256 token")
	}
	if status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}
}

func TestValidateJWT_RejectsWrongSecret(t *testing.T) {
	token, err := auth.GenerateJWT("a-different-secret", 1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	reached, status := serveWithAuth(t, "Bearer "+token)

	if reached {
		t.Error("handler should not be reached with a token signed by another secret")
	}
	if status != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, status)
	}
}
