package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Role names - kept identical to the backend so pages and API calls agree.
const (
	RoleAdmin   = "ADMIN"
	RoleManager = "MANAGER"
	RoleUser    = "USER"
)

// RequireRoles renders a 403 page instead of a JSON error (this is a web
// app, not an API), for any signed-in user whose role isn't allowed.
// Must run after JWTAuth.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextRoleKey)
		if !exists {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok || !allowed[role] {
			c.HTML(http.StatusForbidden, "layout/forbidden.html", gin.H{
				"Title": "Access Denied",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
