package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	api *client.APIClient
}

func NewCustomerHandler(api *client.APIClient) *CustomerHandler {
	return &CustomerHandler{api: api}
}

func (h *CustomerHandler) List(c *gin.Context) {
	var customers []model.Customer
	data := baseData(c, "Customers")
	if err := h.api.Get("/customers", token(c), &customers); err != nil {
		data["Error"] = "Failed to load customers: " + err.Error()
	}
	data["Customers"] = customers
	c.HTML(http.StatusOK, "customer_template/page_view_customer.html", data)
}

func (h *CustomerHandler) Detail(c *gin.Context) {
	path := "/customers/" + c.Param("custNumber") + "/" + c.Param("custID")

	var customer model.Customer
	if err := h.api.Get(path, token(c), &customer); err != nil {
		c.Redirect(http.StatusFound, "/customers?err="+url.QueryEscape("Customer not found"))
		return
	}

	data := baseData(c, "Customer Detail")
	data["Customer"] = customer
	c.HTML(http.StatusOK, "customer_template/detailscustomer.html", data)
}

func (h *CustomerHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Customer")
	data["IsEdit"] = false
	data["Customer"] = model.Customer{}
	c.HTML(http.StatusOK, "customer_template/create_update.html", data)
}

func (h *CustomerHandler) formBody(c *gin.Context) gin.H {
	return gin.H{
		"cust_name":      c.PostForm("cust_name"),
		"address":        c.PostForm("address"),
		"city":           c.PostForm("city"),
		"contact_person": c.PostForm("contact_person"),
		"email":          c.PostForm("email"),
		"bank_number":    c.PostForm("bank_number"),
		"phone_number":   c.PostForm("phone_number"),
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	custNumber, _ := strconv.Atoi(c.PostForm("cust_number"))
	body := h.formBody(c)
	body["cust_number"] = custNumber
	body["cust_id"] = c.PostForm("cust_id")

	if err := h.api.Post("/customers", token(c), body, nil); err != nil {
		data := baseData(c, "New Customer")
		data["IsEdit"] = false
		data["Customer"] = model.Customer{CustNumber: custNumber, CustID: c.PostForm("cust_id"), CustName: c.PostForm("cust_name")}
		data["Error"] = "Failed to create customer: " + err.Error()
		c.HTML(http.StatusOK, "customer_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/customers?ok="+url.QueryEscape("Customer created"))
}

func (h *CustomerHandler) ShowEdit(c *gin.Context) {
	path := "/customers/" + c.Param("custNumber") + "/" + c.Param("custID")

	var customer model.Customer
	if err := h.api.Get(path, token(c), &customer); err != nil {
		c.Redirect(http.StatusFound, "/customers?err="+url.QueryEscape("Customer not found"))
		return
	}

	data := baseData(c, "Edit Customer")
	data["IsEdit"] = true
	data["Customer"] = customer
	c.HTML(http.StatusOK, "customer_template/create_update.html", data)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	custNumber := c.Param("custNumber")
	custID := c.Param("custID")
	path := "/customers/" + custNumber + "/" + custID

	body := h.formBody(c)

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		cn, _ := strconv.Atoi(custNumber)
		data := baseData(c, "Edit Customer")
		data["IsEdit"] = true
		data["Customer"] = model.Customer{CustNumber: cn, CustID: custID, CustName: c.PostForm("cust_name")}
		data["Error"] = "Failed to update customer: " + err.Error()
		c.HTML(http.StatusOK, "customer_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/customers?ok="+url.QueryEscape("Customer updated"))
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	path := "/customers/" + c.Param("custNumber") + "/" + c.Param("custID")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/customers?err="+url.QueryEscape("Failed to delete customer: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/customers?ok="+url.QueryEscape("Customer deleted"))
}
