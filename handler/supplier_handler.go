package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	api *client.APIClient
}

func NewSupplierHandler(api *client.APIClient) *SupplierHandler {
	return &SupplierHandler{api: api}
}

func (h *SupplierHandler) List(c *gin.Context) {
	var suppliers []model.Supplier
	data := baseData(c, "Suppliers")
	if err := h.api.Get("/suppliers", token(c), &suppliers); err != nil {
		data["Error"] = "Failed to load suppliers: " + err.Error()
	}
	data["Suppliers"] = suppliers
	c.HTML(http.StatusOK, "supplier_template/page_view_supplier.html", data)
}

func (h *SupplierHandler) Detail(c *gin.Context) {
	path := "/suppliers/" + c.Param("supplierNumber") + "/" + c.Param("supplierID")

	var supplier model.Supplier
	if err := h.api.Get(path, token(c), &supplier); err != nil {
		c.Redirect(http.StatusFound, "/suppliers?err="+url.QueryEscape("Supplier not found"))
		return
	}

	data := baseData(c, "Supplier Detail")
	data["Supplier"] = supplier
	c.HTML(http.StatusOK, "supplier_template/detailssupplier.html", data)
}

func (h *SupplierHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Supplier")
	data["IsEdit"] = false
	data["Supplier"] = model.Supplier{}
	c.HTML(http.StatusOK, "supplier_template/create_update.html", data)
}

func (h *SupplierHandler) formBody(c *gin.Context) gin.H {
	return gin.H{
		"supplier_name":  c.PostForm("supplier_name"),
		"address":        c.PostForm("address"),
		"city":           c.PostForm("city"),
		"contact_person": c.PostForm("contact_person"),
		"email":          c.PostForm("email"),
		"phone_number":   c.PostForm("phone_number"),
		"bank_number":    c.PostForm("bank_number"),
	}
}

func (h *SupplierHandler) Create(c *gin.Context) {
	body := h.formBody(c)
	body["supplier_id"] = c.PostForm("supplier_id")

	if err := h.api.Post("/suppliers", token(c), body, nil); err != nil {
		data := baseData(c, "New Supplier")
		data["IsEdit"] = false
		data["Supplier"] = model.Supplier{SupplierID: c.PostForm("supplier_id"), SupplierName: c.PostForm("supplier_name")}
		data["Error"] = "Failed to create supplier: " + err.Error()
		c.HTML(http.StatusOK, "supplier_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/suppliers?ok="+url.QueryEscape("Supplier created"))
}

func (h *SupplierHandler) ShowEdit(c *gin.Context) {
	path := "/suppliers/" + c.Param("supplierNumber") + "/" + c.Param("supplierID")

	var supplier model.Supplier
	if err := h.api.Get(path, token(c), &supplier); err != nil {
		c.Redirect(http.StatusFound, "/suppliers?err="+url.QueryEscape("Supplier not found"))
		return
	}

	data := baseData(c, "Edit Supplier")
	data["IsEdit"] = true
	data["Supplier"] = supplier
	c.HTML(http.StatusOK, "supplier_template/create_update.html", data)
}

func (h *SupplierHandler) Update(c *gin.Context) {
	supplierNumber := c.Param("supplierNumber")
	supplierID := c.Param("supplierID")
	path := "/suppliers/" + supplierNumber + "/" + supplierID

	body := h.formBody(c)

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		sn, _ := strconv.Atoi(supplierNumber)
		data := baseData(c, "Edit Supplier")
		data["IsEdit"] = true
		data["Supplier"] = model.Supplier{SupplierNumber: sn, SupplierID: supplierID, SupplierName: c.PostForm("supplier_name")}
		data["Error"] = "Failed to update supplier: " + err.Error()
		c.HTML(http.StatusOK, "supplier_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/suppliers?ok="+url.QueryEscape("Supplier updated"))
}

func (h *SupplierHandler) Delete(c *gin.Context) {
	path := "/suppliers/" + c.Param("supplierNumber") + "/" + c.Param("supplierID")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/suppliers?err="+url.QueryEscape("Failed to delete supplier: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/suppliers?ok="+url.QueryEscape("Supplier deleted"))
}
