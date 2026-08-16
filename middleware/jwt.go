package middleware

import (
	"net/http"

	"webApps_order-management/auth"

	"github.com/gin-gonic/gin"
)

const (
	ContextTokenKey    = "token"
	ContextUserIDKey   = "userID"
	ContextUsernameKey = "username"
	ContextRoleKey     = "role"
)

// JWTAuth reads the session cookie, validates it, and injects the user's
// identity into the Gin context. Unlike the backend API (which returns a
// 401 JSON body), this is a web app: a missing/expired token redirects the
// browser to the login page instead.
func JWTAuth(secret, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(cookieName)
		if err != nil || tokenString == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(secret, tokenString)
		if err != nil {
			c.SetCookie(cookieName, "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set(ContextTokenKey, tokenString)
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextRoleKey, claims.RoleName)

		c.Next()
	}
}

// OptionalAuth is like JWTAuth but never redirects - used on public pages
// (like the login page itself) that still want to know if someone is
// already signed in.
func OptionalAuth(secret, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(cookieName)
		if err == nil && tokenString != "" {
			if claims, err := auth.ParseToken(secret, tokenString); err == nil {
				c.Set(ContextTokenKey, tokenString)
				c.Set(ContextUserIDKey, claims.UserID)
				c.Set(ContextUsernameKey, claims.Username)
				c.Set(ContextRoleKey, claims.RoleName)
			}
		}
		c.Next()
	}
}
