package handler

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	api *client.APIClient
}

func NewOrderHandler(api *client.APIClient) *OrderHandler {
	return &OrderHandler{api: api}
}

// loadLookups fetches suppliers, customers, and products so the order form
// can offer autocomplete-by-name while still submitting the underlying ID
// (SupplierID / CustID / ProductID) that the backend expects.
func (h *OrderHandler) loadLookups(c *gin.Context, data gin.H) {
	var suppliers []model.Supplier
	if err := h.api.Get("/suppliers", token(c), &suppliers); err == nil {
		data["Suppliers"] = suppliers
	}
	var customers []model.Customer
	if err := h.api.Get("/customers", token(c), &customers); err == nil {
		data["Customers"] = customers
	}
	var products []model.Product
	if err := h.api.Get("/products", token(c), &products); err == nil {
		data["Products"] = products
	}
}

func (h *OrderHandler) List(c *gin.Context) {
	var orders []model.OrderMaster
	data := baseData(c, "Orders")
	if err := h.api.Get("/orders", token(c), &orders); err != nil {
		data["Error"] = "Failed to load orders: " + err.Error()
	}
	data["Orders"] = orders
	c.HTML(http.StatusOK, "order_template/page_view_order.html", data)
}

func (h *OrderHandler) Detail(c *gin.Context) {
	path := "/orders/" + c.Param("orderID") + "/" + c.Param("orderNo")

	var order model.OrderMaster
	if err := h.api.Get(path, token(c), &order); err != nil {
		c.Redirect(http.StatusFound, "/orders?err="+url.QueryEscape("Order not found"))
		return
	}

	data := baseData(c, "Order Detail")
	data["Order"] = order
	var products []model.Product
	if err := h.api.Get("/products", token(c), &products); err == nil {
		data["Products"] = products
	}
	c.HTML(http.StatusOK, "order_template/detailsorder.html", data)
}

func (h *OrderHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Order")
	data["IsEdit"] = false
	data["Order"] = model.OrderMaster{}
	h.loadLookups(c, data)
	c.HTML(http.StatusOK, "order_template/create_update.html", data)
}

func parseFormDate(value string) interface{} {
	if value == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", value); err == nil {
		return t
	}
	return nil
}

// Create -> POST /orders/create : master fields + at least one detail line.
func (h *OrderHandler) Create(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.PostForm("order_id"))
	orderNo := c.PostForm("order_no")

	itemNames := c.PostFormArray("item_name[]")
	measures := c.PostFormArray("measure[]")
	qtys := c.PostFormArray("qty[]")
	prices := c.PostFormArray("price[]")

	details := make([]gin.H, 0, len(itemNames))
	for i := range itemNames {
		if itemNames[i] == "" {
			continue
		}
		qty, _ := strconv.Atoi(qtys[i])
		price, _ := strconv.ParseFloat(prices[i], 64)
		measure := ""
		if i < len(measures) {
			measure = measures[i]
		}
		details = append(details, gin.H{
			"item_name": itemNames[i],
			"measure":   measure,
			"qty":       qty,
			"price":     price,
		})
	}

	body := gin.H{
		"order_id":        orderID,
		"order_no":        orderNo,
		"supplier_from":   c.PostForm("supplier_from"),
		"customer_id":     c.PostForm("customer_id"),
		"description":     c.PostForm("description"),
		"document_number": c.PostForm("document_number"),
		"details":         details,
	}
	if d := parseFormDate(c.PostForm("order_date")); d != nil {
		body["order_date"] = d
	}

	if err := h.api.Post("/orders", token(c), body, nil); err != nil {
		data := baseData(c, "New Order")
		data["IsEdit"] = false
		data["Order"] = model.OrderMaster{OrderID: orderID, OrderNo: orderNo}
		data["Error"] = "Failed to create order: " + err.Error()
		h.loadLookups(c, data)
		c.HTML(http.StatusOK, "order_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/orders?ok="+url.QueryEscape("Order created"))
}

func (h *OrderHandler) ShowEdit(c *gin.Context) {
	path := "/orders/" + c.Param("orderID") + "/" + c.Param("orderNo")

	var order model.OrderMaster
	if err := h.api.Get(path, token(c), &order); err != nil {
		c.Redirect(http.StatusFound, "/orders?err="+url.QueryEscape("Order not found"))
		return
	}

	data := baseData(c, "Edit Order")
	data["IsEdit"] = true
	data["Order"] = order
	h.loadLookups(c, data)
	c.HTML(http.StatusOK, "order_template/create_update.html", data)
}

// Update -> POST /orders/:orderID/:orderNo/edit : master fields only.
// Detail lines are managed separately from the order detail page.
func (h *OrderHandler) Update(c *gin.Context) {
	orderID := c.Param("orderID")
	orderNo := c.Param("orderNo")
	path := "/orders/" + orderID + "/" + orderNo

	body := gin.H{
		"supplier_from":   c.PostForm("supplier_from"),
		"customer_id":     c.PostForm("customer_id"),
		"description":     c.PostForm("description"),
		"document_number": c.PostForm("document_number"),
	}
	if d := parseFormDate(c.PostForm("order_date")); d != nil {
		body["order_date"] = d
	}

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		oid, _ := strconv.Atoi(orderID)
		data := baseData(c, "Edit Order")
		data["IsEdit"] = true
		data["Order"] = model.OrderMaster{OrderID: oid, OrderNo: orderNo}
		data["Error"] = "Failed to update order: " + err.Error()
		h.loadLookups(c, data)
		c.HTML(http.StatusOK, "order_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/orders?ok="+url.QueryEscape("Order updated"))
}

func (h *OrderHandler) Delete(c *gin.Context) {
	path := "/orders/" + c.Param("orderID") + "/" + c.Param("orderNo")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/orders?err="+url.QueryEscape("Failed to delete order: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/orders?ok="+url.QueryEscape("Order deleted"))
}

// --- Detail line endpoints ---

func (h *OrderHandler) AddDetail(c *gin.Context) {
	orderID := c.Param("orderID")
	orderNo := c.Param("orderNo")

	qty, _ := strconv.Atoi(c.PostForm("qty"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

	body := gin.H{
		"item_name": c.PostForm("item_name"),
		"measure":   c.PostForm("measure"),
		"qty":       qty,
		"price":     price,
	}

	redirectBase := "/orders/" + orderID + "/" + orderNo

	if err := h.api.Post(redirectBase+"/details", token(c), body, nil); err != nil {
		c.Redirect(http.StatusFound, redirectBase+"?err="+url.QueryEscape("Failed to add detail: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, redirectBase+"?ok="+url.QueryEscape("Detail line added"))
}

func (h *OrderHandler) DeleteDetail(c *gin.Context) {
	orderID := c.Param("orderID")
	orderNo := c.Param("orderNo")
	orderDetailNo := c.Param("orderDetailNo")

	redirectBase := "/orders/" + orderID + "/" + orderNo
	path := redirectBase + "/details/" + orderDetailNo

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, redirectBase+"?err="+url.QueryEscape("Failed to remove detail: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, redirectBase+"?ok="+url.QueryEscape("Detail line removed"))
}
