package models

type OrderDetail struct {
	ID      uint `json:"id" gorm:"primarykey"`
	Qty     uint
	Amount  uint
	BookID  uint
	Book    Book
	OrderID uint
}
