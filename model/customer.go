package model

type Customer struct {
	CustNumber    int    `json:"cust_number"`
	CustID        string `json:"cust_id"`
	CustName      string `json:"cust_name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email"`
	BankNumber    string `json:"bank_number"`
	PhoneNumber   string `json:"phone_number"`
}
