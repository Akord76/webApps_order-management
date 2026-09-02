package model

type CommitmentFee struct {
	CommitmentID       int `json:"commitment_id"`
	EmployeeCardNumber int `json:"employee_card_number"` //di ambil dari employee (Autocomplete)
	CompleteName 	   string  `json:"complete_name"`
	CommitmentValue 	float64 `json:"commitment_value"`
	ParameterFee    	int     `json:"parameter_fee"`
	ProductID       	string  `json:"product_id"` //diambil dari table product
	ProductName  	   string `json:"product_name"`
}

// Struct untuk mendaratkan response dari API backend
type EmployeeSearchResponse struct {
	Data []struct {
		EmployeeID         int    `json:"employee_id"`
		EmployeeCardNumber int    `json:"employee_card_number"`
		FirstName          string `json:"first_name"`
		LastName           string `json:"last_name"`
		Email              string `json:"email"`
	} `json:"data"`
	Success bool `json:"success"`
}

type ProductSearchResponse struct {
	Data []struct {
		ProductNumber     int    `json:"product_number"`
		ProductID         string `json:"product_id"`
		ProductCategoryID *int   `json:"product_category_id"`
		ProductName       string `json:"product_name"`
		Measure           string `json:"measure"`
		Description       string `json:"description"`
		DocumentNumber    string `json:"document_number"`
	} `json:"data"`
	Success bool `json:"success"`
}
