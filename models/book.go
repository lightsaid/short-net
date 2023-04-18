package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null"`
	Price     uint           `json:"price" gorm:"int;not null"`
	Picture   string         `json:"picture" gorm:"type:varchar(255);not null"`
	Stock     uint           `json:"stock" gorm:"int;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
