package model

type Company struct {
	CompanyNumber int    `json:"company_number"`
	CompanyID     string `json:"company_id"`
	CompanyName   string `json:"company_name"`
	Address       string `json:"address"`
	ContactPerson string `json:"contact_person"`
}
