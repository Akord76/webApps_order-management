package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	api *client.APIClient
}

func NewCategoryHandler(api *client.APIClient) *CategoryHandler {
	return &CategoryHandler{api: api}
}

// List -> GET /categories
func (h *CategoryHandler) List(c *gin.Context) {
	var categories []model.Category
	if err := h.api.Get("/categories", token(c), &categories); err != nil {
		data := baseData(c, "Categories")
		data["Error"] = "Failed to load categories: " + err.Error()
		c.HTML(http.StatusOK, "category_template/page_view_category.html", data)
		return
	}

	data := baseData(c, "Categories")
	data["Categories"] = categories
	c.HTML(http.StatusOK, "category_template/page_view_category.html", data)
}

// Detail -> GET /categories/:id
func (h *CategoryHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	var category model.Category
	if err := h.api.Get("/categories/"+id, token(c), &category); err != nil {
		c.Redirect(http.StatusFound, "/categories?err="+url.QueryEscape("Category not found"))
		return
	}

	data := baseData(c, "Category Detail")
	data["Category"] = category
	c.HTML(http.StatusOK, "category_template/detailscategory.html", data)
}

// ShowCreate -> GET /categories/create
func (h *CategoryHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Category")
	data["IsEdit"] = false
	data["Category"] = model.Category{}
	c.HTML(http.StatusOK, "category_template/create_update.html", data)
}

// Create -> POST /categories/create
func (h *CategoryHandler) Create(c *gin.Context) {
	body := gin.H{"category_name": c.PostForm("category_name")}

	if err := h.api.Post("/categories", token(c), body, nil); err != nil {
		data := baseData(c, "New Category")
		data["IsEdit"] = false
		data["Category"] = model.Category{CategoryName: c.PostForm("category_name")}
		data["Error"] = "Failed to create category: " + err.Error()
		c.HTML(http.StatusOK, "category_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/categories?ok="+url.QueryEscape("Category created"))
}

// ShowEdit -> GET /categories/:id/edit
func (h *CategoryHandler) ShowEdit(c *gin.Context) {
	id := c.Param("id")

	var category model.Category
	if err := h.api.Get("/categories/"+id, token(c), &category); err != nil {
		c.Redirect(http.StatusFound, "/categories?err="+url.QueryEscape("Category not found"))
		return
	}

	data := baseData(c, "Edit Category")
	data["IsEdit"] = true
	data["Category"] = category
	c.HTML(http.StatusOK, "category_template/create_update.html", data)
}

// Update -> POST /categories/:id/edit
func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	body := gin.H{"category_name": c.PostForm("category_name")}

	if err := h.api.Put("/categories/"+id, token(c), body, nil); err != nil {
		data := baseData(c, "Edit Category")
		data["IsEdit"] = true
		data["Category"] = model.Category{CategoryID: idInt, CategoryName: c.PostForm("category_name")}
		data["Error"] = "Failed to update category: " + err.Error()
		c.HTML(http.StatusOK, "category_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/categories?ok="+url.QueryEscape("Category updated"))
}

// Delete -> POST /categories/:id/delete
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.api.Delete("/categories/"+id, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/categories?err="+url.QueryEscape("Failed to delete category: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/categories?ok="+url.QueryEscape("Category deleted"))
}
