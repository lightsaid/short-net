package dbrepo

import "github.com/lightsaid/short-net/models"

func (r *repository) CreateOrderDetail(detail *models.OrderDetail) error {
	return r.DB.Create(detail).Error
}
