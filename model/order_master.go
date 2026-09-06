package model

import "time"

type OrderMaster struct {
	OrderID        int        `json:"order_id"`
	OrderNo        string     `json:"order_no"`
	OrderDate      *time.Time `json:"order_date"`
	SupplierFrom   string     `json:"supplier_from"`
	CustomerID     string     `json:"customer_id"`
	CustomerName   string     `json:"customer_name"`
	Description    string     `json:"description"`
	DocumentNumber string     `json:"document_number"`
	StatusPO       int     `json:"status_po"`
	Status_Po       string     `json:"statuspo"`

	Details []OrderDetail `json:"details,omitempty"`
}
