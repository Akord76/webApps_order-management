package model

import "time"

// Struct Response / Data Utama
type SharingProfit struct {
	SharingProfitNum   int       `json:"sharing_profit_num"`
	SharingProfitID    string    `json:"sharing_profit_id"`
	SharingProfitDate  string    `json:"sharing_profit_date"`
	CommitmentID       int       `json:"commitment_id"`
	EmployeeCardNumber int       `json:"employee_card_number"`
	EmployeeName       string    `json:"employee_name"`
	OrderDoDetailNo    string    `json:"order_do_dtno"`
	OrderNo            string    `json:"order_no"`
	ProductID          string    `json:"product_id"`
	ItemName           string    `json:"item_name"`
	Measure            string    `json:"measure"`
	Qty                int       `json:"qty"`
	ShareValue         float64   `json:"share_value"`
	CreatedAt          time.Time `json:"created_at"`
}

// Struct Request Binding Form POST/PUT ke API Backend
type CreateSharingProfitRequest struct {
	SharingProfitNum  int     `form:"sharing_profit_num" json:"sharing_profit_num" binding:"required"`
	SharingProfitID   string  `form:"sharing_profit_id" json:"sharing_profit_id" binding:"required"`
	SharingProfitDate string  `form:"sharing_profit_date" json:"sharing_profit_date" binding:"required"`
	CommitmentID      int     `form:"commitment_id" json:"commitment_id" binding:"required"`
	OrderDoDetailNo   string  `form:"order_do_dtno" json:"order_do_dtno" binding:"required"`
	Qty               int     `form:"qty" json:"qty" binding:"required,min=1"`
	ShareValue        float64 `form:"share_value" json:"share_value" binding:"required"`
}