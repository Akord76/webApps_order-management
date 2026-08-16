package model

import "time"

type OrderSaleMaster struct {
	OrderSaleNo   int               `json:"order_sale_no"`
	OrderSaleID   string            `json:"order_sale_id"`
	OrderSaleDate *time.Time        `json:"order_sale_date"`
	CustomerID    string            `json:"customer_id"`
	Shipment      string            `json:"shipment"`
	ShipNumber    string            `json:"ship_number"`
	DriverNumber  string            `json:"driver_number"`
	Description   string            `json:"description"`
	Details       []OrderSaleDetail `json:"details,omitempty"`
}
