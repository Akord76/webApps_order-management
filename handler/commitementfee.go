package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type CommitmentFeeHandler struct {
	api *client.APIClient
}

func NewCommitmentFeeHandler(api *client.APIClient) *CommitmentFeeHandler {
	return &CommitmentFeeHandler{api: api}
}

func (h *CommitmentFeeHandler) List(c *gin.Context) {
	// Struct untuk menampung elemen individual commitment fee
	type CommitmentFeeItem struct {
		model.CommitmentFee
		CompleteName string `json:"complete_name"`
		ProductName  string `json:"product_name"`
	}

	// Struct pembungkus sesuai bentuk response JSON dari API
	type APIResponse struct {
		Data    []CommitmentFeeItem `json:"data"` // Key "data" menyesuaikan response API
		Success bool                `json:"success"`
	}

	var apiRes APIResponse
	data := baseData(c, "Commitment Fee")

	// Panggil API ke struct pembungkus (apiRes)
	if err := h.api.Get("/commitmentfees", token(c), &apiRes); err != nil {
		data["Error"] = "Failed to load Commitment Fee records: " + err.Error()
	}

	// Kirim array di dalam apiRes.Data ke template
	data["CommitmentFee"] = apiRes.Data
	c.HTML(http.StatusOK, "commitmentfee_template/page_view_commitmentfee.html", data)
}

func (h *CommitmentFeeHandler) Detail(c *gin.Context) {
	path := "/commitmentfees/" + "/" + c.Param("id")

	var sp model.CommitmentFee
	if err := h.api.Get(path, token(c), &sp); err != nil {
		c.Redirect(http.StatusFound, "/commitmentfees?err="+url.QueryEscape("Commitment Fee record not found"))
		return
	}

	data := baseData(c, "Commitment Fee Detail")
	data["CommitmentFee"] = sp
	c.HTML(http.StatusOK, "commitmentfee_template/details_commitmentfee.html", data)
}

func (h *CommitmentFeeHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Commitment Fee")
	data["IsEdit"] = false
	data["CommitmentFee"] = model.CommitmentFee{}
	c.HTML(http.StatusOK, "commitmentfee_template/create_update.html", data)
}

func (h *CommitmentFeeHandler) formBody(c *gin.Context) map[string]interface{} {
	empCard, _ := strconv.Atoi(c.PostForm("employee_card_number"))
	commVal, _ := strconv.ParseFloat(c.PostForm("commitment_value"), 64)
	paramFee, _ := strconv.Atoi(c.PostForm("parameter_fee"))

	// WAJIB: Ambil product_id dari PostForm
	productID := c.PostForm("product_id")

	return gin.H{
		"employee_card_number": empCard,
		"commitment_value":     commVal,
		"parameter_fee":        paramFee,
		"product_id":           productID, // Tambahkan key ini agar tidak kosong
	}
}

func (h *CommitmentFeeHandler) Create(c *gin.Context) {
	// 1. Ambil seluruh data dari form input HTML
	empCard, _ := strconv.Atoi(c.PostForm("employee_card_number"))
	commVal, _ := strconv.ParseFloat(c.PostForm("commitment_value"), 64)
	paramFee, _ := strconv.Atoi(c.PostForm("parameter_fee"))
	prodID := c.PostForm("product_id") // Ambil ProductID dari hidden/text input HTML

	// 2. Buat payload request ke API Backend
	body := gin.H{
		"employee_card_number": empCard,
		"commitment_value":     commVal,
		"parameter_fee":        paramFee,
		"product_id":           prodID, // Kirim product_id ke API
	}

	// 3. Eksekusi request POST ke backend microservice
	if err := h.api.Post("/commitmentfees", token(c), body, nil); err != nil {
		data := baseData(c, "New Commitment Fee")
		data["IsEdit"] = false

		// Re-populate nilai ke model agar input form tetap terisi saat render halaman error
		data["CommitmentFee"] = model.CommitmentFee{
			EmployeeCardNumber: empCard,
			CommitmentValue:    commVal,
			ParameterFee:       paramFee,
			ProductID:          prodID,
		}
		data["Error"] = "Failed to create Commitment Fee record: " + err.Error()
		c.HTML(http.StatusOK, "commitmentfee_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/commitmentfees?ok="+url.QueryEscape("Commitment Fee record created"))
}

func (h *CommitmentFeeHandler) ShowEdit(c *gin.Context) {
	path := "/commitmentfees/" + c.Param("num")

	var sp model.CommitmentFee
	if err := h.api.Get(path, token(c), &sp); err != nil {
		c.Redirect(http.StatusFound, "/commitmentfees?err="+url.QueryEscape("Commitment Fee record not found"))
		return
	}

	data := baseData(c, "Edit Commitment Fee")
	data["IsEdit"] = true
	data["CommitmentFee"] = sp
	c.HTML(http.StatusOK, "commitmentfee_template/create_update.html", data)
}

func (h *CommitmentFeeHandler) Update(c *gin.Context) {
	num := c.Param("num")
	path := "/commitmentfees/" + num

	body := h.formBody(c)

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		comID, _ := strconv.Atoi(num)
		data := baseData(c, "Edit Commitment Fee")
		data["IsEdit"] = true
		data["CommitmentFee"] = model.CommitmentFee{CommitmentID: comID}
		data["Error"] = "Failed to update Commitment Fee record: " + err.Error()
		c.HTML(http.StatusOK, "commitmentfee_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/commitmentfees?ok="+url.QueryEscape("Commitment Fee record updated"))
}

func (h *CommitmentFeeHandler) Delete(c *gin.Context) {
	path := "/commitmentfees/" + c.Param("num")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/commitmentfees?err="+url.QueryEscape("Failed to delete Commitment Fee record: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/commitmentfees?ok="+url.QueryEscape("Commitment Fee record deleted"))
}

// Handler untuk endpoint Autocomplete JSON
func (h *CommitmentFeeHandler) EmployeeAutocomplete(c *gin.Context) {
	searchQuery := c.Query("q")
	if searchQuery == "" {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
		return
	}

	// Buat URL endpoint API Backend
	apiURL := fmt.Sprintf("/commitmentfees/employees/search?search=%s", searchQuery)

	var apiRes model.EmployeeSearchResponse
	if err := h.api.Get(apiURL, token(c), &apiRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan data JSON ke Frontend
	c.JSON(http.StatusOK, apiRes)
}

func (h *CommitmentFeeHandler) ProductAutocomplete(c *gin.Context) {
	searchQuery := c.Query("q")
	if searchQuery == "" {
		c.JSON(http.StatusOK, gin.H{"results": []string{}})
		return
	}

	apiURL := fmt.Sprintf("/commitmentfees/products/search?search=%s", searchQuery)

	var apiRes model.ProductSearchResponse
	if err := h.api.Get(apiURL, token(c), &apiRes); err != nil {
		// LOG ERROR KE TERMINAL BACKEND (Penting untuk melihat error aslinya)
		log.Printf("[ProductAutocomplete ERROR]: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	// Cetak log respon untuk memastikan isi struct tidak kosong
	log.Printf("[ProductAutocomplete SUCCESS]: %+v", apiRes)

	c.JSON(http.StatusOK, apiRes)
}
