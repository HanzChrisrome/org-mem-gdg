package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Clean up visitors before each test
	mu.Lock()
	visitors = make(map[string]*visitor)
	mu.Unlock()

	t.Run("allows requests within limit", func(t *testing.T) {
		r := gin.New()
		// Low limit for testing: 100 RPS, 1 Burst
		r.Use(RateLimiter(100, 1))
		r.GET("/test-allowed", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test-allowed", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("blocks requests exceeding burst", func(t *testing.T) {
		mu.Lock()
		visitors = make(map[string]*visitor)
		mu.Unlock()

		r := gin.New()
		// Limit: 1 RPS, 1 Burst
		r.Use(RateLimiter(1, 1))
		r.GET("/test-burst", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// First request allowed
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test-burst", nil)
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second request blocked immediately
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test-burst", nil)
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	})

    t.Run("allows requests after cooldown", func(t *testing.T) {
		mu.Lock()
		visitors = make(map[string]*visitor)
		mu.Unlock()

		r := gin.New()
		// Limit: 10 RPS (100ms interval), 1 Burst
		r.Use(RateLimiter(10, 1))
		r.GET("/test-cooldown", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// First request allowed
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test-cooldown", nil)
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Wait for token to replenish (10 RPS = 1 token per 100ms)
		time.Sleep(150 * time.Millisecond)

		// Second request should be allowed now
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test-cooldown", nil)
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)
	})
}
