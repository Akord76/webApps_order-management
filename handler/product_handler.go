package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	api *client.APIClient
}

func NewProductHandler(api *client.APIClient) *ProductHandler {
	return &ProductHandler{api: api}
}

func (h *ProductHandler) List(c *gin.Context) {
	var products []model.Product
	data := baseData(c, "Products")
	if err := h.api.Get("/products", token(c), &products); err != nil {
		data["Error"] = "Failed to load products: " + err.Error()
	}
	data["Products"] = products
	c.HTML(http.StatusOK, "product_template/page_view_product.html", data)
}

func (h *ProductHandler) Detail(c *gin.Context) {
	path := "/products/" + c.Param("productNumber") + "/" + c.Param("productID")

	var product model.Product
	if err := h.api.Get(path, token(c), &product); err != nil {
		c.Redirect(http.StatusFound, "/products?err="+url.QueryEscape("Product not found"))
		return
	}

	data := baseData(c, "Product Detail")
	data["Product"] = product
	c.HTML(http.StatusOK, "product_template/detailsproduct.html", data)
}

func (h *ProductHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Product")
	data["IsEdit"] = false
	data["Product"] = model.Product{}
	c.HTML(http.StatusOK, "product_template/create_update.html", data)
}

func (h *ProductHandler) productBodyFromForm(c *gin.Context) gin.H {
	body := gin.H{
		"product_name":    c.PostForm("product_name"),
		"measure":         c.PostForm("measure"),
		"description":     c.PostForm("description"),
		"document_number": c.PostForm("document_number"),
	}
	if v := c.PostForm("product_category_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			body["product_category_id"] = n
		}
	}
	return body
}

func (h *ProductHandler) Create(c *gin.Context) {
	productNumber, _ := strconv.Atoi(c.PostForm("product_number"))
	body := h.productBodyFromForm(c)
	body["product_number"] = productNumber
	body["product_id"] = c.PostForm("product_id")

	if err := h.api.Post("/products", token(c), body, nil); err != nil {
		data := baseData(c, "New Product")
		data["IsEdit"] = false
		data["Product"] = model.Product{
			ProductNumber: productNumber, ProductID: c.PostForm("product_id"),
			ProductName: c.PostForm("product_name"), Measure: c.PostForm("measure"),
			Description: c.PostForm("description"), DocumentNumber: c.PostForm("document_number"),
		}
		data["Error"] = "Failed to create product: " + err.Error()
		c.HTML(http.StatusOK, "product_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/products?ok="+url.QueryEscape("Product created"))
}

func (h *ProductHandler) ShowEdit(c *gin.Context) {
	path := "/products/" + c.Param("productNumber") + "/" + c.Param("productID")

	var product model.Product
	if err := h.api.Get(path, token(c), &product); err != nil {
		c.Redirect(http.StatusFound, "/products?err="+url.QueryEscape("Product not found"))
		return
	}

	data := baseData(c, "Edit Product")
	data["IsEdit"] = true
	data["Product"] = product
	c.HTML(http.StatusOK, "product_template/create_update.html", data)
}

func (h *ProductHandler) Update(c *gin.Context) {
	productNumber := c.Param("productNumber")
	productID := c.Param("productID")
	path := "/products/" + productNumber + "/" + productID

	body := h.productBodyFromForm(c)

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		pn, _ := strconv.Atoi(productNumber)
		data := baseData(c, "Edit Product")
		data["IsEdit"] = true
		data["Product"] = model.Product{
			ProductNumber: pn, ProductID: productID,
			ProductName: c.PostForm("product_name"), Measure: c.PostForm("measure"),
			Description: c.PostForm("description"), DocumentNumber: c.PostForm("document_number"),
		}
		data["Error"] = "Failed to update product: " + err.Error()
		c.HTML(http.StatusOK, "product_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/products?ok="+url.QueryEscape("Product updated"))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	path := "/products/" + c.Param("productNumber") + "/" + c.Param("productID")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/products?err="+url.QueryEscape("Failed to delete product: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/products?ok="+url.QueryEscape("Product deleted"))
}
