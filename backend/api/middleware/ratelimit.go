package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	ips sync.Map
}

type IPState struct {
	limiter    *rate.Limiter
	failures   int
	blockUntil time.Time
	mu         sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}
func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)

		stateInterface, _ := rl.ips.LoadOrStore(ip, &IPState{
			limiter: rate.NewLimiter(rate.Every(time.Minute/5), 5),
		})
		state := stateInterface.(*IPState)

		state.mu.Lock()
		defer state.mu.Unlock()

		if time.Now().Before(state.blockUntil) {
			waitDuration := time.Until(state.blockUntil)
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", waitDuration.Seconds()))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Too many failed attempts. Please try again in %s.", waitDuration.Round(time.Second)),
			})
			return
		}

		if !state.limiter.Allow() {
			w.Header().Set("Retry-After", "60")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Rate limit exceeded. Please wait a moment.",
			})
			return
		}

		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		if rec.statusCode == http.StatusUnauthorized {
			state.failures++
			log.Printf("[Warning] Failed login attempt %d from IP: %s", state.failures, ip)

			var blockDuration time.Duration

			if state.failures >= 15 {
				blockDuration = 15 * time.Minute
			} else if state.failures >= 10 {
				blockDuration = 5 * time.Minute
			} else if state.failures >= 5 {
				blockDuration = 1 * time.Minute
			}

			if blockDuration > 0 {
				state.blockUntil = time.Now().Add(blockDuration)
				log.Printf("[Alert] Blocking IP %s for %v due to repeated failures", ip, blockDuration)
			}

		} else if rec.statusCode == http.StatusOK {
			if state.failures > 0 {
				state.failures = 0
				state.blockUntil = time.Time{}
			}
		}
	})
}
