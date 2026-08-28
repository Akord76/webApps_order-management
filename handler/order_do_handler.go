package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type OrderDoHandler struct {
	api *client.APIClient
}

func NewOrderDoHandler(api *client.APIClient) *OrderDoHandler {
	return &OrderDoHandler{api: api}
}

// LookupCustomers proxies GET /api/customers/getCustomers/autocomplete/:custName.
// The browser calls THIS endpoint (same-origin, cookie auth) instead of the
// backend directly - avoids CORS and never exposes the JWT to client-side JS.
func (h *OrderDoHandler) LookupCustomers(c *gin.Context) {
	custName := c.Param("custName")

	var result interface{}
	path := "/customers/getCustomers/autocomplete/" + url.PathEscape(custName)
	if err := h.api.Get(path, token(c), &result); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "lookup failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// LookupOrders proxies GET /api/orders/getOrder/autocomplete/:custID/:query.
func (h *OrderDoHandler) LookupOrders(c *gin.Context) {
	custID := c.Param("custID")
	query := c.Param("query")

	var result interface{}
	path := "/orders/getOrder/autocomplete/" + url.PathEscape(custID) + "/" + url.PathEscape(query)
	if err := h.api.Get(path, token(c), &result); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "lookup failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *OrderDoHandler) loadLookupsDo(c *gin.Context, data gin.H) {
	var customers []model.Customer
	if err := h.api.Get("/customers", token(c), &customers); err == nil {
		data["Customers"] = customers
	}
	var products []model.Product
	if err := h.api.Get("/products", token(c), &products); err == nil {
		data["Products"] = products
	}
}

func (h *OrderDoHandler) List(c *gin.Context) {
	var orderDo []model.OrderDoMaster
	data := baseData(c, "Delivery Order")
	if err := h.api.Get("/orderdo", token(c), &orderDo); err != nil {
		data["Error"] = "Failed to load order sales: " + err.Error()
	}
	data["OrderDo"] = orderDo
	c.HTML(http.StatusOK, "order_do_template/page_view_order_do.html", data)
}

func (h *OrderDoHandler) Detail(c *gin.Context) {
	path := "/orderdo/" + c.Param("orderDoID") + "/" + c.Param("orderDoNo")

	var orderDo model.OrderDoMaster
	if err := h.api.Get(path, token(c), &orderDo); err != nil {
		c.Redirect(http.StatusFound, "/orderdo?err="+url.QueryEscape("Order Do not found"))
		return
	}

	// 2. Hitung Grand Total di handler yang sama
	// var grandTotal float64
	// for _, d := range orderDoDetails {
	// 	grandTotal += float64(d.Qty) * d.Price
	// }

	data := baseData(c, "Order Do Detail")
	data["OrderDo"] = orderDo
	data["GrandTotal"] = calculateGrandTotal(orderDo.Details)
	c.HTML(http.StatusOK, "order_do_template/detailsorder_do.html", data)
}

// Helper function terpisah
func calculateGrandTotal(details []model.OrderDoDetail) float64 {
	var total float64
	for _, d := range details {
		total += float64(d.Qty) * d.Price
	}
	return total
}



func (h *OrderDoHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Order Sale")
	data["IsEdit"] = false
	data["OrderDo"] = model.OrderDoMaster{}
	h.loadLookupsDo(c, data)
	c.HTML(http.StatusOK, "order_do_template/create_update.html", data)
}

func (h *OrderDoHandler) Create(c *gin.Context) {
	orderDoID, _ := strconv.Atoi(c.PostForm("order_do_id"))
	orderDoNo := c.PostForm("order_do_no")

	detailIDs := c.PostFormArray("order_do_dtid[]")
	if len(detailIDs) == 0 {
		detailIDs = c.PostFormArray("order_do_dtno[]")
	}
	orderNos := c.PostFormArray("order_no[]")
	productIDs := c.PostFormArray("product_id[]")
	itemNames := c.PostFormArray("item_name[]")
	measures := c.PostFormArray("measure[]")
	qtys := c.PostFormArray("qty[]")
	prices := c.PostFormArray("price[]")
	documentNumbers := c.PostFormArray("document_number[]")
	details := make([]gin.H, 0, len(productIDs))
	for i := range productIDs {
		if productIDs[i] == "" {
			continue
		}

		qty, _ := strconv.Atoi(qtys[i])
		price, _ := strconv.ParseFloat(prices[i], 64)

		detailID := ""
		if i < len(detailIDs) {
			detailID = detailIDs[i]
		}
		orderNo := ""
		if i < len(orderNos) {
			orderNo = orderNos[i]
		}
		itemName := ""
		if i < len(itemNames) {
			itemName = itemNames[i]
		}
		measure := ""
		if i < len(measures) {
			measure = measures[i]
		}

		documentNumber := ""
		if i < len(documentNumbers) {
			documentNumber = documentNumbers[i]
		}
		details = append(details, gin.H{
			"order_do_dtno":   detailID,
			"order_no":        orderNo,
			"product_id":      productIDs[i], // String product_id dari array
			"item_name":       itemName,
			"measure":         measure,
			"qty":             qty,
			"price":           price,
			"document_number": documentNumber,
		})
	}

	body := gin.H{
		"order_do_id":   orderDoID,
		"order_do_no":   orderDoNo,
		"customer_id":   c.PostForm("customer_id"),
		"shipment":      c.PostForm("shipment"),
		"ship_number":   c.PostForm("ship_number"),
		"driver_number": c.PostForm("driver_number"),
		"description":   c.PostForm("description"),
		"details":       details,
	}

	if d := parseFormDate(c.PostForm("order_do_date")); d != nil {
		body["order_do_date"] = d
	}

	if err := h.api.Post("/orderdo", token(c), body, nil); err != nil {
		data := baseData(c, "New Order Do")
		data["IsEdit"] = false
		data["OrderDo"] = model.OrderDoMaster{OrderDoID: orderDoID, OrderDoNo: orderDoNo}
		data["Error"] = "Failed to create order sale: " + err.Error()
		h.loadLookupsDo(c, data)
		c.HTML(http.StatusOK, "order_do_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/orderdo?ok="+url.QueryEscape("Order sale created"))
}

func (h *OrderDoHandler) ShowEdit(c *gin.Context) {
	path := "/orderdo/" + c.Param("orderDoID") + "/" + c.Param("orderDoNo")

	var orderDoMaster model.OrderDoMaster
	if err := h.api.Get(path, token(c), &orderDoMaster); err != nil {
		c.Redirect(http.StatusFound, "/orderdo?err="+url.QueryEscape("Order DO not found"))
		return
	}

	data := baseData(c, "Edit Order DO")
	data["IsEdit"] = true
	data["OrderDo"] = orderDoMaster
	h.loadLookupsDo(c, data)
	c.HTML(http.StatusOK, "order_do_template/create_update.html", data)
}


// Update -> POST /orderdo/:orderDoNo/:orderDoID/edit : master fields only.
func (h *OrderDoHandler) Update(c *gin.Context) {

	orderDoID := c.Param("orderDoID")
	orderDoNo := c.Param("orderDoNo")
	path := "/orderdo/" + orderDoID + "/" + orderDoNo

	body := gin.H{
		"customer_id":   c.PostForm("customer_id"),
		"shipment":      c.PostForm("shipment"),
		"ship_number":   c.PostForm("ship_number"),
		"driver_number": c.PostForm("driver_number"),
		"description":   c.PostForm("description"),
	}
	if d := parseFormDate(c.PostForm("order_do_date")); d != nil {
		body["order_do_date"] = d
	}

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		osn, _ := strconv.Atoi(orderDoID)
		data := baseData(c, "Edit Order Sale")
		data["IsEdit"] = true
		data["OrderDo"] = model.OrderDoMaster{OrderDoID: osn, OrderDoNo: orderDoNo}
		data["Error"] = "Failed to update order sale: " + err.Error()
		h.loadLookupsDo(c, data)
		c.HTML(http.StatusOK, "order_do_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/orderdo?ok="+url.QueryEscape("Order sale updated"))
}

func (h *OrderDoHandler) Delete(c *gin.Context) {
	path := "/orderdo/" + c.Param("orderDoID") + "/" + c.Param("orderDoNo")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/orderdo?err="+url.QueryEscape("Failed to delete order sale: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/orderdo?ok="+url.QueryEscape("Order sale deleted"))
}

// --- Detail line endpoints ---

func (h *OrderDoHandler) AddDetail(c *gin.Context) {
	orderDoID := c.Param("orderDoID")
	orderDoNo := c.Param("orderDoNo")

	qty, _ := strconv.Atoi(c.PostForm("qty"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

	body := gin.H{
		"order_do_dtno": c.PostForm("order_do_dtno"),
		"order_no":      c.PostForm("order_no"),
		"product_id":    c.PostForm("product_id"),
		"item_name":     c.PostForm("item_name"),
		"measure":       c.PostForm("measure"),
		"qty":           qty,
		"price":         price,
	}

	// Alamat endpoint API Backend
	apiPath := "/orderdo/" + orderDoNo + "/details"

	if err := h.api.Post(apiPath, token(c), body, nil); err != nil {
		// Jika request menggunakan AJAX/Fetch
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" || c.ContentType() == "application/json" {
			c.String(http.StatusBadRequest, "Gagal menambah detail: "+err.Error())
			return
		}
		// Jika request menggunakan Form submit biasa
		c.Redirect(http.StatusFound, "/orderdo/"+orderDoID+"/"+orderDoNo+"?err="+url.QueryEscape(err.Error()))
		return
	}

	// Respon sukses: Kembalikan 201 Created jika AJAX, atau Redirect jika HTML Form
	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.Status(http.StatusCreated) // HTTP 201
		return
	}

	c.Redirect(http.StatusFound, "/orderdo/"+orderDoID+"/"+orderDoNo+"?ok="+url.QueryEscape("Detail line added"))
}

// DeleteDetail removes a single detail line directly - confirmed via
// Postman: DELETE /api/orderdo/details/:orderDoDetailID/:orderDoDetailNo.

func (h *OrderDoHandler) DeleteDetail(c *gin.Context) {
	targetDetailID := c.Param("orderDoDetailID")
	targetDetailNo := c.Param("orderDoDetailNo")

	// Panggil backend API menggunakan method Delete
	apiPath := "/orderdo/details/" + targetDetailID + "/" + targetDetailNo
	if err := h.api.Delete(apiPath, token(c)); err != nil {
		// Kirim error HTTP 500 agar ditangkap blok catch(error) di JS
		c.String(http.StatusInternalServerError, "Failed to remove detail: "+err.Error())
		return
	}

	// Kirim respon HTTP 204 agar ditangkap blok try di JS
	c.String(http.StatusNoContent, "Detail line removed")
}
