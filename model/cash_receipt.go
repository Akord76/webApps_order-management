package model

type CashReceipt struct {
	ChasNumber int     `json:"chas_number"`
	Name       string  `json:"name"`
	ChasValue  float64 `json:"chas_value"`
	Notes      string  `json:"notes"`
}
