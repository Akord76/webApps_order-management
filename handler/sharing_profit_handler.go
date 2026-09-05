package handler

import (
	"fmt"
	"net/http"
	"net/url"
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

// List menampilkan daftar Sharing Profit
func (h *SharingProfitHandler) List(c *gin.Context) {
	var list []model.SharingProfit
	data := baseData(c, "Sharing Profit")

	if err := h.api.Get("/sharing-profits", token(c), &list); err != nil {
		data["Error"] = "Failed to load sharing profit records: " + err.Error()
	}
	data["SharingProfits"] = list

	c.HTML(http.StatusOK, "sharing_profit_template/page_view_sharing_profit.html", data)
}

// Detail menampilkan detail 1 record Sharing Profit
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

// ShowCreateForm menampilkan halaman form create Sharing Profit
func (h *SharingProfitHandler) ShowCreateForm(c *gin.Context) {
	orderDoID := c.Query("order_do_id")
	orderDoNo := c.Query("order_do_dtno")

	data := baseData(c, "New Sharing Profit")
	var sp model.SharingProfit

	// Default date hari ini
	sp.SharingProfitDate = time.Now()
	sp.Qty = 1

	// Pre-fill jika dibuka dari Detail DO
	if orderDoID != "" && orderDoNo != "" {
		apiURL := fmt.Sprintf("/load_details/%s/%s/details", orderDoID, url.PathEscape(orderDoNo))

		var apiRes struct {
			Data    []model.OrderDoDetail `json:"data"`
			Success bool                  `json:"success"`
		}

		if err := h.api.Get(apiURL, token(c), &apiRes); err == nil && apiRes.Success && len(apiRes.Data) > 0 {
			detail := apiRes.Data[0]
			sp.OrderDoNo = detail.OrderDoNo
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

	c.HTML(http.StatusOK, "sharing_profit_template/create_update.html", data)
}

// ProcessBulk memproses otomasisasi kalkulasi Sharing Profit dari DO
func (h *SharingProfitHandler) ProcessBulk(c *gin.Context) {
	orderDONo := c.PostForm("order_do_dtno")
	if orderDONo == "" {
		c.Redirect(http.StatusSeeOther, "/sharing-profits?err="+url.QueryEscape("Nomor DO wajib diisi"))
		return
	}

	payload := gin.H{
		"order_do_dtno": orderDONo,
	}

	var apiRes struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err := h.api.Post("/sharing-profits/process-bulk", token(c), payload, &apiRes); err != nil {
		c.Redirect(http.StatusSeeOther, "/sharing-profits?err="+url.QueryEscape("Bulk sharing profit Failed: "+err.Error()))
		return
	}

	c.Redirect(http.StatusSeeOther, "/sharing-profits?ok="+url.QueryEscape("Sharing profit berhasil diproses dari DO "+orderDONo))
}

// Create menyimpan record Sharing Profit manual/baru
func (h *SharingProfitHandler) Create(c *gin.Context) {
	var req model.CreateSharingProfitRequest
	if err := c.ShouldBind(&req); err != nil {
		data := baseData(c, "New Sharing Profit")
		data["IsEdit"] = false
		data["Error"] = "Form input tidak valid: " + err.Error()
		c.HTML(http.StatusBadRequest, "sharing_profit_template/create_update.html", data)
		return
	}

	var apiRes struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err := h.api.Post("/sharing-profits", token(c), req, &apiRes); err != nil {
		data := baseData(c, "New Sharing Profit")
		data["IsEdit"] = false
		data["Error"] = "Failed to create sharing profit record: " + err.Error()
		c.HTML(http.StatusOK, "sharing_profit_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/sharing-profits?ok="+url.QueryEscape("Sharing profit record created"))
}

// ShowEditForm menampilkan halaman form edit Sharing Profit dengan data eksisting
func (h *SharingProfitHandler) ShowEditForm(c *gin.Context) {
	num := c.Param("num")
	id := c.Param("id")

	data := baseData(c, "Edit Sharing Profit")

	apiURL := fmt.Sprintf("/sharing-profits/detail/%s/%s", num, id)
	var apiRes struct {
		Data    model.SharingProfit `json:"data"`
		Success bool                `json:"success"`
		Message string              `json:"message"`
	}

	if err := h.api.Get(apiURL, token(c), &apiRes); err != nil || !apiRes.Success {
		c.Redirect(http.StatusSeeOther, "/sharing-profits?err="+url.QueryEscape("Gagal mengambil data Sharing Profit eksisting"))
		return
	}

	// Fetch data pendukung DO jika OrderDoNo tersedia
	if apiRes.Data.OrderDoNo != "" {
		doApiURL := fmt.Sprintf("/orderdo/detail-by-dono/%s", apiRes.Data.OrderDoNo)
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

	c.HTML(http.StatusOK, "sharing_profit_template/create_update.html", data)
}

// Update memproses pembaharuan data Sharing Profit
func (h *SharingProfitHandler) Update(c *gin.Context) {
	num := c.Param("num")
	id := c.Param("id")

	var req model.CreateSharingProfitRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "sharing_profit_template/create_update.html", gin.H{
			"Error":         "Form input tidak valid: " + err.Error(),
			"IsEdit":        true,
			"SharingProfit": req,
		})
		return
	}

	apiURL := fmt.Sprintf("/sharing-profits/%s/%s/update", num, id)
	var apiRes struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err := h.api.Put(apiURL, token(c), req, &apiRes); err != nil || !apiRes.Success {
		c.HTML(http.StatusInternalServerError, "sharing_profit_template/create_update.html", gin.H{
			"Error":         "Gagal memperbarui Sharing Profit: " + apiRes.Message,
			"IsEdit":        true,
			"SharingProfit": req,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/sharing-profits?ok="+url.QueryEscape("Sharing Profit berhasil diperbarui"))
}

// Delete menghapus data Sharing Profit
func (h *SharingProfitHandler) Delete(c *gin.Context) {
	path := "/sharing-profits/" + c.Param("num") + "/" + c.Param("id")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/sharing-profits?err="+url.QueryEscape("Failed to delete sharing profit record: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/sharing-profits?ok="+url.QueryEscape("Sharing profit record deleted"))
}

// Handler untuk endpoint Autocomplete JSON
func (h *SharingProfitHandler) EmployeeAutocomplete(c *gin.Context) {
	searchQuery := c.Query("q")
	if searchQuery == "" {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
		return
	}

	// Buat URL endpoint API Backend
	apiURL := fmt.Sprintf("/sharing-profits/employees/search?search=%s", searchQuery)

	var apiRes model.EmployeeSearchResponse
	if err := h.api.Get(apiURL, token(c), &apiRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan data JSON ke Frontend
	c.JSON(http.StatusOK, apiRes)
}

