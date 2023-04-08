package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"type:varchar(32);not null;"`
	Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null;"`
	Password  string         `json:"-" gorm:"type:varchar(64);not null;"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255);not null;"`
	Role      string         `json:"role" gorm:"type:varchar(64);default:'USER';not null;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 一对多关系
	Links []Link `json:"links" gorm:"constraint:OnDelete:CASCADE;"`
}
