package model

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type OrderDoDetail struct {
	OrderDoDetailID int     `json:"order_do_dtid"`
	OrderDoDetailNo string  `json:"order_do_dtno"`
	OrderNo         string  `json:"order_no"`
	ProductID       string  `json:"product_id"`
	ItemName        string  `json:"item_name"` //di ambil dari master product join product
	Measure         string  `json:"measure"`   //di ambil dari master product join product
	Qty             int     `json:"qty"`
	Price           float64 `json:"price"`

	DocumentNumber string `json:"document_number"`
}

// Di file model struct Go Anda
func (d OrderDoDetail) Total() string {
	tot := float64(d.Qty) * d.Price
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", tot) // Menghasilkan misal: "Rp 500.000"
}

func (d OrderDoDetail) FormatPrice() string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", d.Price)
}

func (d OrderDoDetail) FormatTotal() string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", float64(d.Qty)*d.Price)
}

func (o OrderDoMaster) GrandTotal() string {
	var total float64
	for _, d := range o.Details {
		total += float64(d.Qty) * d.Price
	}
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", total)
}
