package model

type OrderDoDetail struct {
	OrderDoDetailID int     `json:"order_do_dtid"`
	OrderDoNo       string  `json:"order_do_dtno"`
	OrderNo         string  `json:"order_no"`
	ItemName        string  `json:"item_name"`
	Measure         string  `json:"measure"`
	Qty             int     `json:"qty"`
	Price           float64 `json:"price"`
	DocumentNumber  string  `json:"document_number"`
}
