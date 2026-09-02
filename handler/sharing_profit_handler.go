package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

// ShowEditForm menampilkan halaman form edit Sharing Profit dengan data eksisting dari API Backend
func (h *SharingProfitHandler) ShowEditForm(c *gin.Context) {
	num := c.Param("num")
	id := c.Param("id")

	data := baseData(c, "Edit Sharing Profit")

	// 1. Tembak API Backend untuk mengambil data Sharing Profit eksisting
	apiURL := fmt.Sprintf("/sharing-profits/detail/%s/%s", num, id)
	var apiRes struct {
		Data    model.SharingProfit `json:"data"`
		Success bool                `json:"success"`
		Message string              `json:"message"`
	}

	if err := h.api.Get(apiURL, token(c), &apiRes); err != nil || !apiRes.Success {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Gagal mengambil data Sharing Profit eksisting.",
		})
		return
	}

	// 2. Jika data DO detail ada, ambil detail DO tambahan (Product Name, Measure, Order No) untuk tampilan
	if apiRes.Data.OrderDoDetailNo != "" {
		doApiURL := fmt.Sprintf("/orderdo/detail-by-dtno/%s", apiRes.Data.OrderDoDetailNo)
		var doRes struct {
			Data struct {
				OrderNo  string `json:"order_no"`
				ItemName string `json:"item_name"`
				Measure  string `json:"measure"`
			} `json:"data"`
			Success bool `json:"success"`
		}
		if err := h.api.Get(doApiURL, token(c), &doRes); err == nil && doRes.Success {
			apiRes.Data.OrderNo = doRes.Data.OrderNo
			apiRes.Data.ItemName = doRes.Data.ItemName
			apiRes.Data.Measure = doRes.Data.Measure
		}
	}

	data["SharingProfit"] = apiRes.Data
	data["IsEdit"] = true

	c.HTML(http.StatusOK, "sharing_profit_create_update.html", data)
}

// Update memproses request form POST/PUT untuk memperbarui data Sharing Profit
// Update memproses request form POST/PUT untuk memperbarui data Sharing Profit
func (h *SharingProfitHandler) Update(c *gin.Context) {
	num := c.Param("num")
	id := c.Param("id")

	var req model.CreateSharingProfitRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "sharing_profit_create_update.html", gin.H{
			"Error":         "Form input tidak valid: " + err.Error(),
			"IsEdit":        true,
			"SharingProfit": req,
		})
		return
	}

	// Tembak API Backend untuk memperbarui data
	apiURL := fmt.Sprintf("/sharing-profits/%s/%s/update", num, id)
	var apiRes struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	// PERBAIKAN: Sesuaikan urutan parameter (URL, token, payload request, target response)
	if err := h.api.Put(apiURL, token(c), req, &apiRes); err != nil || !apiRes.Success {
		c.HTML(http.StatusInternalServerError, "sharing_profit_create_update.html", gin.H{
			"Error":         "Gagal memperbarui Sharing Profit: " + apiRes.Message,
			"IsEdit":        true,
			"SharingProfit": req,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/sharing-profits?success=Sharing+Profit+berhasil+diperbarui")
}

func (h *SharingProfitHandler) Delete(c *gin.Context) {
	path := "/sharing-profits/" + c.Param("num") + "/" + c.Param("id")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/sharing-profits?err="+url.QueryEscape("Failed to delete sharing profit record: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/sharing-profits?ok="+url.QueryEscape("Sharing profit record deleted"))
}

func (h *SharingProfitHandler) EmployeeAutocomplete(c *gin.Context) {
	searchQuery := c.Param("employeeCardNumber")
	if searchQuery == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []model.CommitmentFee{},
		})
		return
	}

	// Endpoint API Backend
	apiURL := "/commitmentfees/employees/" + searchQuery

	// Wrapper struct sesuai bentuk JSON API
	var apiRes struct {
		Data    []model.CommitmentFee `json:"data"`
		Success bool                  `json:"success"`
	}

	if err := h.api.Get(apiURL, token(c), &apiRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Kirimkan respon JSON balik ke frontend
	c.JSON(http.StatusOK, gin.H{
		"success": apiRes.Success,
		"data":    apiRes.Data,
	})
}

// ShowCreateForm menampilkan halaman form create Sharing Profit dengan data DO pre-filled

func (h *SharingProfitHandler) ShowCreateForm(c *gin.Context) {
	orderDoID := c.Query("order_do_id")
	orderDoNo := c.Query("order_do_no")

	data := baseData(c, "New Sharing Profit")
	var sp model.SharingProfit

	// Set tanggal default ke hari ini (YYYY-MM-DD)
	sp.SharingProfitDate = time.Now().Format("2006-01-02")

	// Jika diakses dari Halaman Detail DO (Opsi 3)
	if orderDoID != "" && orderDoNo != "" {
		apiURL := fmt.Sprintf("/load_details/%s/%s/details", orderDoID, url.PathEscape(orderDoNo))

		var apiRes struct {
			Data    []model.OrderDoDetail `json:"data"`
			Success bool                  `json:"success"`
		}

		if err := h.api.Get(apiURL, token(c), &apiRes); err == nil && apiRes.Success && len(apiRes.Data) > 0 {
			detail := apiRes.Data[0]
			sp.OrderDoDetailNo = detail.OrderDoDetailNo
			sp.OrderNo = detail.OrderNo
			sp.ProductID = detail.ProductID
			sp.ItemName = detail.ItemName
			sp.Measure = detail.Measure
			sp.Qty = detail.Qty
		} else {
			data["Warning"] = "Data Detail DO tidak ditemukan atau gagal dimuat."
		}
	}

	data["SharingProfit"] = sp
	data["IsEdit"] = false

	c.HTML(http.StatusOK, "sharing_profit_create_update.html", data)
}

