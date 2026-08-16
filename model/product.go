package model

type Product struct {
	ProductNumber     int       `json:"product_number"`
	ProductID         string    `json:"product_id"`
	ProductCategoryID *int      `json:"product_category_id"`
	ProductName       string    `json:"product_name"`
	Measure           string    `json:"measure"`
	Description       string    `json:"description"`
	DocumentNumber    string    `json:"document_number"`
	Category          *Category `json:"category,omitempty"`
}
