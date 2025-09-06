package middleware

import (
	"github.com/go-dev-frame/sponge/pkg/rails"

	"github.com/gin-gonic/gin"
)

// RailsCookieAuthMiddleware validates and decrypts a Rails encrypted cookie,
// attaches the session payload to context under key "rails_session".
func RailsCookieAuthMiddleware(secretKeyBase string, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing cookie"})
			return
		}

		session, err := rails.DecodeSignedCookie(secretKeyBase, cookie, cookieName)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid cookie"})
			return
		}

		c.Set("rails_session", session)
		c.Next()
	}
}
