package model

type OrderSaleDetail struct {
	OrderSaleDetailNo int     `json:"order_sale_detail_no"`
	OrderSaleDetailID string  `json:"order_sale_detail_id"`
	OrderSaleNo       int     `json:"order_sale_no"`
	ItemName          string  `json:"item_name"`
	Measure           string  `json:"measure"`
	Qty               int     `json:"qty"`
	Price             float64 `json:"price"`
	DocumentNumber    string  `json:"document_number"`
}
