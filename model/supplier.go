package model

type Supplier struct {
	SupplierNumber int    `json:"supplier_number"`
	SupplierID     string `json:"supplier_id"`
	SupplierName   string `json:"supplier_name"`
	Address        string `json:"address"`
	City           string `json:"city"`
	ContactPerson  string `json:"contact_person"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	BankNumber     string `json:"bank_number"`
}
