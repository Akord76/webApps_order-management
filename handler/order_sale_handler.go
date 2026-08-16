package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type OrderSaleHandler struct {
	api *client.APIClient
}

func NewOrderSaleHandler(api *client.APIClient) *OrderSaleHandler {
	return &OrderSaleHandler{api: api}
}

func (h *OrderSaleHandler) List(c *gin.Context) {
	var orderSales []model.OrderSaleMaster
	data := baseData(c, "Order Sales")
	if err := h.api.Get("/order-sales", token(c), &orderSales); err != nil {
		data["Error"] = "Failed to load order sales: " + err.Error()
	}
	data["OrderSales"] = orderSales
	c.HTML(http.StatusOK, "order_sale_template/page_view_order_sale.html", data)
}

func (h *OrderSaleHandler) Detail(c *gin.Context) {
	path := "/order-sales/" + c.Param("orderSaleNo") + "/" + c.Param("orderSaleID")

	var orderSale model.OrderSaleMaster
	if err := h.api.Get(path, token(c), &orderSale); err != nil {
		c.Redirect(http.StatusFound, "/order-sales?err="+url.QueryEscape("Order sale not found"))
		return
	}

	data := baseData(c, "Order Sale Detail")
	data["OrderSale"] = orderSale
	c.HTML(http.StatusOK, "order_sale_template/detailsorder_sale.html", data)
}

func (h *OrderSaleHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Order Sale")
	data["IsEdit"] = false
	data["OrderSale"] = model.OrderSaleMaster{}
	c.HTML(http.StatusOK, "order_sale_template/create_update.html", data)
}

// Create -> POST /order-sales/create : master fields + at least one detail line.
func (h *OrderSaleHandler) Create(c *gin.Context) {
	orderSaleNo, _ := strconv.Atoi(c.PostForm("order_sale_no"))
	orderSaleID := c.PostForm("order_sale_id")

	detailIDs := c.PostFormArray("order_sale_detail_id[]")
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
		detailID := ""
		if i < len(detailIDs) {
			detailID = detailIDs[i]
		}
		details = append(details, gin.H{
			"order_sale_detail_id": detailID,
			"item_name":            itemNames[i],
			"measure":              measure,
			"qty":                  qty,
			"price":                price,
		})
	}

	body := gin.H{
		"order_sale_no": orderSaleNo,
		"order_sale_id": orderSaleID,
		"customer_id":   c.PostForm("customer_id"),
		"shipment":      c.PostForm("shipment"),
		"ship_number":   c.PostForm("ship_number"),
		"driver_number": c.PostForm("driver_number"),
		"description":   c.PostForm("description"),
		"details":       details,
	}
	if d := parseFormDate(c.PostForm("order_sale_date")); d != nil {
		body["order_sale_date"] = d
	}

	if err := h.api.Post("/order-sales", token(c), body, nil); err != nil {
		data := baseData(c, "New Order Sale")
		data["IsEdit"] = false
		data["OrderSale"] = model.OrderSaleMaster{OrderSaleNo: orderSaleNo, OrderSaleID: orderSaleID}
		data["Error"] = "Failed to create order sale: " + err.Error()
		c.HTML(http.StatusOK, "order_sale_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/order-sales?ok="+url.QueryEscape("Order sale created"))
}

func (h *OrderSaleHandler) ShowEdit(c *gin.Context) {
	path := "/order-sales/" + c.Param("orderSaleNo") + "/" + c.Param("orderSaleID")

	var orderSale model.OrderSaleMaster
	if err := h.api.Get(path, token(c), &orderSale); err != nil {
		c.Redirect(http.StatusFound, "/order-sales?err="+url.QueryEscape("Order sale not found"))
		return
	}

	data := baseData(c, "Edit Order Sale")
	data["IsEdit"] = true
	data["OrderSale"] = orderSale
	c.HTML(http.StatusOK, "order_sale_template/create_update.html", data)
}

// Update -> POST /order-sales/:orderSaleNo/:orderSaleID/edit : master fields only.
func (h *OrderSaleHandler) Update(c *gin.Context) {
	orderSaleNo := c.Param("orderSaleNo")
	orderSaleID := c.Param("orderSaleID")
	path := "/order-sales/" + orderSaleNo + "/" + orderSaleID

	body := gin.H{
		"customer_id":   c.PostForm("customer_id"),
		"shipment":      c.PostForm("shipment"),
		"ship_number":   c.PostForm("ship_number"),
		"driver_number": c.PostForm("driver_number"),
		"description":   c.PostForm("description"),
	}
	if d := parseFormDate(c.PostForm("order_sale_date")); d != nil {
		body["order_sale_date"] = d
	}

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		osn, _ := strconv.Atoi(orderSaleNo)
		data := baseData(c, "Edit Order Sale")
		data["IsEdit"] = true
		data["OrderSale"] = model.OrderSaleMaster{OrderSaleNo: osn, OrderSaleID: orderSaleID}
		data["Error"] = "Failed to update order sale: " + err.Error()
		c.HTML(http.StatusOK, "order_sale_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/order-sales?ok="+url.QueryEscape("Order sale updated"))
}

func (h *OrderSaleHandler) Delete(c *gin.Context) {
	path := "/order-sales/" + c.Param("orderSaleNo") + "/" + c.Param("orderSaleID")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/order-sales?err="+url.QueryEscape("Failed to delete order sale: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/order-sales?ok="+url.QueryEscape("Order sale deleted"))
}

// --- Detail line endpoints ---

func (h *OrderSaleHandler) AddDetail(c *gin.Context) {
	orderSaleNo := c.Param("orderSaleNo")
	orderSaleID := c.Param("orderSaleID")

	qty, _ := strconv.Atoi(c.PostForm("qty"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

	body := gin.H{
		"order_sale_detail_id": c.PostForm("order_sale_detail_id"),
		"item_name":            c.PostForm("item_name"),
		"measure":              c.PostForm("measure"),
		"qty":                  qty,
		"price":                price,
	}

	redirectBase := "/order-sales/" + orderSaleNo + "/" + orderSaleID

	if err := h.api.Post(redirectBase+"/details", token(c), body, nil); err != nil {
		c.Redirect(http.StatusFound, redirectBase+"?err="+url.QueryEscape("Failed to add detail: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, redirectBase+"?ok="+url.QueryEscape("Detail line added"))
}

func (h *OrderSaleHandler) DeleteDetail(c *gin.Context) {
	orderSaleNo := c.Param("orderSaleNo")
	orderSaleID := c.Param("orderSaleID")
	orderSaleDetailNo := c.Param("orderSaleDetailNo")
	orderSaleDetailID := c.Param("orderSaleDetailID")

	redirectBase := "/order-sales/" + orderSaleNo + "/" + orderSaleID
	path := redirectBase + "/details/" + orderSaleDetailNo + "/" + orderSaleDetailID

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, redirectBase+"?err="+url.QueryEscape("Failed to remove detail: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, redirectBase+"?ok="+url.QueryEscape("Detail line removed"))
}
