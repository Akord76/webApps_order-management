package model

import "time"

type SharingProfit struct {
	SharingProfitNum         int        `json:"sharing_profit_num"`
	SharingProfitID          string     `json:"sharing_profit_id"`
	SharingProfitDate        time.Time  `json:"sharing_profit_date"`
	OrderDoNo                string     `json:"order_do_dtno"`
	ProductID                string     `json:"product_id"`
	Qty                      int        `json:"qty"`
	ShareValue               float64    `json:"share_value"`
	TotalShareValue          float64    `json:"total_share_value"`
	CommitmentID             int        `json:"commitment_id"`
	IsSharingProfitProcessed bool       `json:"is_sharing_profit_processed"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty"`

	// Field pendukung untuk JOIN tampilan (menggunakan int sesuai tipe kolom di DB)
	EmployeeCardNumber int    `json:"employee_card_number"`
	EmployeeName       string `json:"employee_name"`
	OrderNo            string `json:"order_no"`
	ItemName           string `json:"item_name"`
	Measure            string `json:"measure"`
}

// Struct Request Binding Form POST/PUT ke API Backend
type CreateSharingProfitRequest struct {
	SharingProfitNum         int       `form:"sharing_profit_num" json:"sharing_profit_num" binding:"required"`
	SharingProfitID          string    `form:"sharing_profit_id" json:"sharing_profit_id" binding:"required"`
	SharingProfitDate        time.Time `form:"sharing_profit_date" json:"sharing_profit_date" binding:"required"`
	OrderDoNo                string    `form:"order_do_dtno" json:"order_do_dtno" binding:"required"`
	ProductID                string    `form:"product_id" json:"product_id" binding:"required"`
	Qty                      int       `form:"qty" json:"qty" binding:"required,min=1"`
	ShareValue               float64   `form:"share_value" json:"share_value" binding:"required"`
	TotalShareValue          float64   `form:"total_share_value" json:"total_share_value" binding:"required"`
	CommitmentID             int       `form:"commitment_id" json:"commitment_id" binding:"required"`
	IsSharingProfitProcessed bool      `form:"is_sharing_profit_processed" json:"is_sharing_profit_processed"`

	CreatedAt time.Time  `form:"column:CreatedAt;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time `form:"column:UpdatedAt;autoUpdateTime" json:"updated_at,omitempty"`

	// Field pendukung untuk JOIN tampilan (menggunakan int sesuai tipe kolom di DB)
	EmployeeCardNumber int    `form:"employee_card_number" json:"employee_card_number"`
	EmployeeName       string `form:"employee_name" json:"employee_name"`
	OrderNo            string `form:"order_no" json:"order_no"`
	ItemName           string `form:"item_name" json:"item_name"`
	Measure            string `form:"measure" json:"measure,omitempty"`
}
