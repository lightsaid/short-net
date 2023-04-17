package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	TotalAmount  uint           `json:"total_amount"`
	UserID       uint           `json:"user_id" gorm:"index"`
	OrderDetails []OrderDetail  `json:"order_details"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
