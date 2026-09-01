package model

type SharingProfit struct {
	SharingProfitNum   int     `json:"sharing_profit_num"`
	SharingProfitID    string  `json:"sharing_profit_id"`
	EmployeeCardNumber *int    `json:"employee_card_number"`//di ambil dari table commitment Fee (Autocomplete)
	OrderDoNo          string  `json:"order_do_dtno"`//di ambil dari order detail DO No, input manual
	ProductID          string  `json:"product_id"`//di ambil dari order detail DO No, (Autocomplete,parameter order dono)
	Qty                int     `json:"qty"`//di ambil dari order detail DO No, suggestion of order dono
	ShareValue         float64 `json:"share_value"`//di ambil dari commitment fee DO No,
}
