package model

type CommitmentFee struct {
	CommitmentID       int     `json:"commitment_id"`
	EmployeeCardNumber *int    `json:"employee_card_number"`
	CommitmentValue    float64 `json:"commitment_value"`
	ParameterFee       *int    `json:"parameter_fee"`
	ProductID          string  `json:"product_id"`
}
