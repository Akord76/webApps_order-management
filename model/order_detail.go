package model

type OrderDetail struct {
	OrderDetailNo int     `json:"order_detail_no"`
	OrderNo       string  `json:"order_no"`
	ItemName      string  `json:"item_name"`
	Measure       string  `json:"measure"`
	Qty           int     `json:"qty"`
	Price         float64 `json:"price"`
}
