package models

import (
	"time"
)

type Link struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	LongURL   string    `json:"long_url" gorm:"type:varchar(255);not null"`
	ShortHash string    `json:"short_hash" gorm:"type:varchar(16);uniqueIndex;not null;"`
	Click     uint      `json:"click" gorm:"default:0;not null;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiredAt time.Time `json:"expired_at" gorm:"index;not null;"`
}
