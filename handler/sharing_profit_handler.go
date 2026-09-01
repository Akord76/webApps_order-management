package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type SharingProfitHandler struct {
	api *client.APIClient
}

func NewSharingProfitHandler(api *client.APIClient) *SharingProfitHandler {
	return &SharingProfitHandler{api: api}
}

func (h *SharingProfitHandler) List(c *gin.Context) {
	var list []model.SharingProfit
	data := baseData(c, "Sharing Profit")
	if err := h.api.Get("/sharing-profits", token(c), &list); err != nil {
		data["Error"] = "Failed to load sharing profit records: " + err.Error()
	}
	data["SharingProfits"] = list
	c.HTML(http.StatusOK, "sharing_profit_template/page_view_sharing_profit.html", data)
}

func (h *SharingProfitHandler) Detail(c *gin.Context) {
	path := "/sharing-profits/" + c.Param("num") + "/" + c.Param("id")

	var sp model.SharingProfit
	if err := h.api.Get(path, token(c), &sp); err != nil {
		c.Redirect(http.StatusFound, "/sharing-profits?err="+url.QueryEscape("Sharing profit record not found"))
		return
	}

	data := baseData(c, "Sharing Profit Detail")
	data["SharingProfit"] = sp
	c.HTML(http.StatusOK, "sharing_profit_template/details_sharing_profit.html", data)
}

func (h *SharingProfitHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Sharing Profit")
	data["IsEdit"] = false
	data["SharingProfit"] = model.SharingProfit{}
	c.HTML(http.StatusOK, "sharing_profit_template/create_update.html", data)
}

func (h *SharingProfitHandler) formBody(c *gin.Context) gin.H {
	body := gin.H{
		"order_sale_detail_id": c.PostForm("order_sale_detail_id"),
	}
	if v := c.PostForm("commitment_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			body["commitment_id"] = n
		}
	}
	if v, err := strconv.ParseFloat(c.PostForm("share_value"), 64); err == nil {
		body["share_value"] = v
	}
	return body
}

func (h *SharingProfitHandler) Create(c *gin.Context) {
	num, _ := strconv.Atoi(c.PostForm("sharing_profit_num"))
	body := h.formBody(c)
	body["sharing_profit_num"] = num
	body["sharing_profit_id"] = c.PostForm("sharing_profit_id")

	if err := h.api.Post("/sharing-profits", token(c), body, nil); err != nil {
		data := baseData(c, "New Sharing Profit")
		data["IsEdit"] = false
		data["SharingProfit"] = model.SharingProfit{
			SharingProfitNum: num, SharingProfitID: c.PostForm("sharing_profit_id"),
		}
		data["Error"] = "Failed to create sharing profit record: " + err.Error()
		c.HTML(http.StatusOK, "sharing_profit_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/sharing-profits?ok="+url.QueryEscape("Sharing profit record created"))
}

func (h *SharingProfitHandler) ShowEdit(c *gin.Context) {
	path := "/sharing-profits/" + c.Param("num") + "/" + c.Param("id")

	var sp model.SharingProfit
	if err := h.api.Get(path, token(c), &sp); err != nil {
		c.Redirect(http.StatusFound, "/sharing-profits?err="+url.QueryEscape("Sharing profit record not found"))
		return
	}

	data := baseData(c, "Edit Sharing Profit")
	data["IsEdit"] = true
	data["SharingProfit"] = sp
	c.HTML(http.StatusOK, "sharing_profit_template/create_update.html", data)
}

func (h *SharingProfitHandler) Update(c *gin.Context) {
	num := c.Param("num")
	id := c.Param("id")
	path := "/sharing-profits/" + num + "/" + id

	body := h.formBody(c)

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		n, _ := strconv.Atoi(num)
		data := baseData(c, "Edit Sharing Profit")
		data["IsEdit"] = true
		data["SharingProfit"] = model.SharingProfit{SharingProfitNum: n, SharingProfitID: id}
		data["Error"] = "Failed to update sharing profit record: " + err.Error()
		c.HTML(http.StatusOK, "sharing_profit_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/sharing-profits?ok="+url.QueryEscape("Sharing profit record updated"))
}

func (h *SharingProfitHandler) Delete(c *gin.Context) {
	path := "/sharing-profits/" + c.Param("num") + "/" + c.Param("id")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/sharing-profits?err="+url.QueryEscape("Failed to delete sharing profit record: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/sharing-profits?ok="+url.QueryEscape("Sharing profit record deleted"))
}
