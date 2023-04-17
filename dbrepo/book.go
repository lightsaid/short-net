package dbrepo

import (
	"errors"

	"github.com/lightsaid/short-net/models"
)

var ErrUnderstock = errors.New("库存不足")

func (r *repository) CreateBook(book *models.Book) error {
	return r.DB.Create(book).Error
}

func (r *repository) ListBooks(f Filters) ([]*models.Book, error) {
	books := make([]*models.Book, 0, f.Limit())
	err := r.DB.Limit(f.Limit()).Offset(f.Offset()).Find(&books).Error
	return books, err
}

func (r *repository) DeductionStock(id uint, qty uint) error {
	var q models.Book
	err := r.DB.First(&q, id).Error
	if err != nil {
		return err
	}

	if q.Stcok < qty {
		return ErrUnderstock
	}

	q.Stcok -= qty

	return r.DB.Save(&q).Error
}
