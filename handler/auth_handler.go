package handler

import (
	"net/http"
	"net/url"

	"webApps_order-management/auth"
	"webApps_order-management/client"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	api        *client.APIClient
	cookieName string
	secureCookie bool
}

func NewAuthHandler(api *client.APIClient, cookieName string, secureCookie bool) *AuthHandler {
	return &AuthHandler{api: api, cookieName: cookieName, secureCookie: secureCookie}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "auth_template/login.html", baseData(c, "Login"))
}

func (h *AuthHandler) ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "auth_template/register.html", baseData(c, "Register"))
}

type loginAPIResponse struct {
	Token string `json:"token"`
	User  struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
		FullName string `json:"full_name"`
		RoleName string `json:"role_name"`
	} `json:"user"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var resp loginAPIResponse
	body := gin.H{"username": username, "password": password}
	if err := h.api.Post("/auth/login", "", body, &resp); err != nil {
		c.Redirect(http.StatusFound, "/login?err="+url.QueryEscape("Invalid username or password"))
		return
	}

	// 8 hours matches the backend's default JWT expiry closely enough for
	// the cookie lifetime; the JWT itself remains the source of truth for
	// actual expiration.
	c.SetCookie(h.cookieName, resp.Token, 8*3600, "/", "", h.secureCookie, true)
	c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	fullName := c.PostForm("full_name")
	email := c.PostForm("email")

	if err := auth.ValidatePasswordStrength(password); err != nil {
		c.Redirect(http.StatusFound, "/register?err="+url.QueryEscape(err.Error()))
		return
	}

	body := gin.H{
		"username":  username,
		"password":  password,
		"full_name": fullName,
		"email":     email,
	}
	if err := h.api.Post("/auth/register", "", body, nil); err != nil {
		c.Redirect(http.StatusFound, "/register?err="+url.QueryEscape("Registration failed: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/login?ok="+url.QueryEscape("Account created. Please sign in."))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(h.cookieName, "", -1, "/", "", h.secureCookie, true)
	c.Redirect(http.StatusFound, "/login")
}
