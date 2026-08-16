package model

type SharingProfit struct {
	SharingProfitNum  int     `json:"sharing_profit_num"`
	SharingProfitID   string  `json:"sharing_profit_id"`
	CommitmentID      *int    `json:"commitment_id"`
	OrderSaleDetailID string  `json:"order_sale_detail_id"`
	ShareValue        float64 `json:"share_value"`
}
