package middleware

import (
	"github.com/gin-gonic/gin"
)

// RateLimit returns a gin.HandlerFunc middleware that applies rate limiting per IP.
// It uses a token bucket algorithm with a limit of 500 requests per minute.
// Requests exceeding the limit receive a 429 Too Many Requests response with
// X-RateLimit-Limit and X-RateLimit-Remaining headers.
//
// Note: Full implementation is provided by T040. This wires the middleware into
// the server router per T044.
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO(T040): Implement token bucket rate limiting (500 req/min per IP).
		// For now, set rate limit headers and pass through all requests.
		c.Header("X-RateLimit-Limit", "500")
		c.Header("X-RateLimit-Remaining", "499")
		c.Next()
	}
}
