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

	data := baseData(c, "Order Do Detail")
	data["OrderDo"] = orderDo
	c.HTML(http.StatusOK, "order_do_template/detailsorder_do.html", data)
}

func (h *OrderDoHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Order Sale")
	data["IsEdit"] = false
	data["OrderDo"] = model.OrderDoMaster{}
	c.HTML(http.StatusOK, "order_do_template/create_update.html", data)
}

// Create -> POST /orderdo/create : master fields + at least one detail line.
func (h *OrderDoHandler) Create(c *gin.Context) {
	orderDoID, _ := strconv.Atoi(c.PostForm("order_do_id"))
	orderDoNo := c.PostForm("order_do_no")

	detailIDs := c.PostFormArray("order_do_detail_id[]")
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
			"order_do_detail_id": detailID,
			"item_name":          itemNames[i],
			"measure":            measure,
			"qty":                qty,
			"price":              price,
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
		c.HTML(http.StatusOK, "order_do_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/orderdo?ok="+url.QueryEscape("Order sale created"))
}

func (h *OrderDoHandler) ShowEdit(c *gin.Context) {
	path := "/orderdo/" + c.Param("orderDoID") + "/" + c.Param("orderDoNo")

	var orderDoMaster model.OrderDoMaster
	if err := h.api.Get(path, token(c), &orderDoMaster); err != nil {
		c.Redirect(http.StatusFound, "/orderdo?err="+url.QueryEscape("Order sale not found"))
		return
	}

	data := baseData(c, "Edit Order Sale")
	data["IsEdit"] = true
	data["OrderDo"] = orderDoMaster
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
		"order_do_detail_id": c.PostForm("order_do_detail_id"),
		"item_name":          c.PostForm("item_name"),
		"measure":            c.PostForm("measure"),
		"qty":                qty,
		"price":              price,
	}

	redirectBase := "/orderdo/" + orderDoID + "/" + orderDoNo

	if err := h.api.Post(redirectBase+"/details", token(c), body, nil); err != nil {
		c.Redirect(http.StatusFound, redirectBase+"?err="+url.QueryEscape("Failed to add detail: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, redirectBase+"?ok="+url.QueryEscape("Detail line added"))
}

func (h *OrderDoHandler) DeleteDetail(c *gin.Context) {
    orderDoID := c.Param("orderDoID")
    orderDoNo := c.Param("orderDoNo")
    targetDetailID := c.Param("orderDoDetailID") // "103", "104", dst
    pagePath := "/orderdo/" + orderDoID + "/" + orderDoNo

    // 1. Ambil semua detail existing
	var details []model.OrderDoDetail
	detailsPath := "/orderdo/" + orderDoID + "/" + orderDoNo + "/details"
	err := h.api.Get(detailsPath, token(c), &details)
    if err != nil {
        c.Redirect(http.StatusFound, pagePath+"?err="+url.QueryEscape("Gagal ambil data: "+err.Error()))
        return
    }

    // 2. Pisahkan baris yang mau dihapus vs yang dipertahankan
    found := false
    remaining := make([]model.OrderDoDetail, 0, len(details))
    for _, d := range details {
		if strconv.Itoa(d.OrderDoDetailID) == targetDetailID {
            found = true
            continue
        }
        remaining = append(remaining, d)
    }
    if !found {
        c.Redirect(http.StatusFound, pagePath+"?err="+url.QueryEscape("Baris detail tidak ditemukan"))
        return
    }

    // 3. Hapus SEMUA baris via API 2-param yang sudah ada
    apiPath := "/orderdo/details/" + orderDoID + "/" + orderDoNo
    if err := h.api.Delete(apiPath, token(c)); err != nil {
        c.Redirect(http.StatusFound, pagePath+"?err="+url.QueryEscape("Failed to remove detail: "+err.Error()))
        return
    }

    // 4. Re-add baris yang tersisa, dengan rollback kalau gagal di tengah
    var addErr error
    for _, d := range remaining {
		body := gin.H{
			"order_do_detail_id": d.OrderDoDetailID,
			"item_name":          d.ItemName,
			"measure":            d.Measure,
			"qty":                d.Qty,
			"price":              d.Price,
		}
		if err := h.api.Post("/orderdo/"+orderDoID+"/"+orderDoNo+"/details", token(c), body, nil); err != nil {
            addErr = err
            break
        }
    }
    if addErr != nil {
        for _, d := range details {
			body := gin.H{
				"order_do_detail_id": d.OrderDoDetailID,
				"item_name":          d.ItemName,
				"measure":            d.Measure,
				"qty":                d.Qty,
				"price":              d.Price,
			}
			_ = h.api.Post("/orderdo/"+orderDoID+"/"+orderDoNo+"/details", token(c), body, nil) // rollback best-effort
        }
        c.Redirect(http.StatusFound, pagePath+"?err="+url.QueryEscape("Gagal hapus baris, data dikembalikan seperti semula: "+addErr.Error()))
        return
    }

    c.Redirect(http.StatusFound, pagePath+"?ok="+url.QueryEscape("Detail line removed"))
}