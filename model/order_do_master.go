package model

import "time"

type OrderDoMaster struct {
	OrderDoID    int        `json:"order_do_id"`
	OrderDoNo    string     `json:"order_do_no"`
	OrderDoDate  *time.Time `json:"order_do_date"`
	CustomerID   string     `json:"customer_id"`
	Shipment     string     `json:"shipment"`
	ShipNumber   string     `json:"ship_number"`
	DriverNumber string     `json:"driver_number"`
	Description  string     `json:"description"`

	Details []OrderDoDetail `json:"details,omitempty"`
}
