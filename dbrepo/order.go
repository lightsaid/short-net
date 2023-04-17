package dbrepo

import "github.com/lightsaid/short-net/models"

func (r *repository) CreateOrder(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *repository) ListOrders(uid uint, f Filters) ([]*models.Order, error) {
	orders := make([]*models.Order, 0, f.Limit())
	err := r.DB.Preload("OrderDetails").Limit(f.Limit()).Offset(f.Offset()).Find(&orders).Error
	return orders, err
}
