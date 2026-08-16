package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	api *client.APIClient
}

func NewUserHandler(api *client.APIClient) *UserHandler {
	return &UserHandler{api: api}
}

// List -> GET /users
func (h *UserHandler) List(c *gin.Context) {
	var users []model.AppUser
	data := baseData(c, "User Management")
	if err := h.api.Get("/users", token(c), &users); err != nil {
		data["Error"] = "Failed to load users: " + err.Error()
	}
	data["Users"] = users
	c.HTML(http.StatusOK, "user_template/page_view_user.html", data)
}

// Detail -> GET /users/:id
func (h *UserHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	var user model.AppUser
	if err := h.api.Get("/users/"+id, token(c), &user); err != nil {
		c.Redirect(http.StatusFound, "/users?err="+url.QueryEscape("User not found"))
		return
	}

	data := baseData(c, "User Detail")
	data["User"] = user
	c.HTML(http.StatusOK, "user_template/detailsuser.html", data)
}

// ShowCreate -> GET /users/create
func (h *UserHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New User")
	data["IsEdit"] = false
	data["User"] = model.AppUser{IsActive: true, RoleName: "USER"}
	c.HTML(http.StatusOK, "user_template/create_update.html", data)
}

// Create -> POST /users/create
func (h *UserHandler) Create(c *gin.Context) {
	body := gin.H{
		"username":  c.PostForm("username"),
		"password":  c.PostForm("password"),
		"full_name": c.PostForm("full_name"),
		"email":     c.PostForm("email"),
		"role_name": c.PostForm("role_name"),
	}

	if err := h.api.Post("/users", token(c), body, nil); err != nil {
		data := baseData(c, "New User")
		data["IsEdit"] = false
		data["User"] = model.AppUser{
			Username: c.PostForm("username"), FullName: c.PostForm("full_name"),
			Email: c.PostForm("email"), RoleName: c.PostForm("role_name"), IsActive: true,
		}
		data["Error"] = "Failed to create user: " + err.Error()
		c.HTML(http.StatusOK, "user_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/users?ok="+url.QueryEscape("User created"))
}

// ShowEdit -> GET /users/:id/edit
func (h *UserHandler) ShowEdit(c *gin.Context) {
	id := c.Param("id")

	var user model.AppUser
	if err := h.api.Get("/users/"+id, token(c), &user); err != nil {
		c.Redirect(http.StatusFound, "/users?err="+url.QueryEscape("User not found"))
		return
	}

	data := baseData(c, "Edit User")
	data["IsEdit"] = true
	data["User"] = user
	c.HTML(http.StatusOK, "user_template/create_update.html", data)
}

// Update -> POST /users/:id/edit (the API's update endpoint doesn't accept
// a password change - only full_name, email, role_name, is_active)
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	body := gin.H{
		"full_name": c.PostForm("full_name"),
		"email":     c.PostForm("email"),
		"role_name": c.PostForm("role_name"),
		"is_active": c.PostForm("is_active") == "true",
	}

	if err := h.api.Put("/users/"+id, token(c), body, nil); err != nil {
		data := baseData(c, "Edit User")
		data["IsEdit"] = true
		data["User"] = model.AppUser{
			UserID: idInt, FullName: c.PostForm("full_name"),
			Email: c.PostForm("email"), RoleName: c.PostForm("role_name"),
		}
		data["Error"] = "Failed to update user: " + err.Error()
		c.HTML(http.StatusOK, "user_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/users?ok="+url.QueryEscape("User updated"))
}

// Delete -> POST /users/:id/delete
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.api.Delete("/users/"+id, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/users?err="+url.QueryEscape("Failed to delete user: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/users?ok="+url.QueryEscape("User deleted"))
}
