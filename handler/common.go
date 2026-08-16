package handler

import (
	"webApps_order-management/middleware"

	"github.com/gin-gonic/gin"
)

// token reads the caller's JWT (set by middleware.JWTAuth) so handlers can
// forward it as a Bearer token to the backend API.
func token(c *gin.Context) string {
	if v, ok := c.Get(middleware.ContextTokenKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func currentUsername(c *gin.Context) string {
	if v, ok := c.Get(middleware.ContextUsernameKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func currentRole(c *gin.Context) string {
	if v, ok := c.Get(middleware.ContextRoleKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func canWrite(c *gin.Context) bool {
	role := currentRole(c)
	return role == middleware.RoleAdmin || role == middleware.RoleManager
}

// baseData seeds every template render with the fields every page/partial
// expects: Title, the signed-in user's identity, and optional flash
// messages carried in via query string after a redirect (?err=...&ok=...).
func baseData(c *gin.Context, title string) gin.H {
	data := gin.H{
		"Title":    title,
		"Username": currentUsername(c),
		"Role":     currentRole(c),
		"CanWrite": canWrite(c),
	}
	if errMsg := c.Query("err"); errMsg != "" {
		data["Error"] = errMsg
	}
	if okMsg := c.Query("ok"); okMsg != "" {
		data["Success"] = okMsg
	}
	return data
}
