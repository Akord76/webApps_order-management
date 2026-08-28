package model
import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type OrderDetail struct {
	OrderDetailNo int     `json:"order_detail_no"`
	OrderNo       string  `json:"order_no"`
	ItemName      string  `json:"item_name"`
	Measure       string  `json:"measure"`
	Qty           int     `json:"qty"`
	Price         float64 `json:"price"`
}


// Di file model struct Go Anda
func (d OrderDetail) Total() string {
	tot := float64(d.Qty) * d.Price
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", tot) // Menghasilkan misal: "Rp 500.000"
}

func (d OrderDetail) FormatPrice() string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", d.Price)
}

func (d OrderDetail) FormatTotal() string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", float64(d.Qty)*d.Price)
}

func (o OrderMaster) GrandTotal() string {
	var total float64
	for _, d := range o.Details {
		total += float64(d.Qty) * d.Price
	}
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.0f", total)
}