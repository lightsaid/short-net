package dbrepo

import (
	"errors"

	"github.com/lightsaid/short-net/models"
)

var ErrUnderstock = errors.New("库存不足")

func (r *repository) CreateBook(book *models.Book) error {
	return r.DB.Create(book).Error
}

func (r *repository) GetBook(id uint) (models.Book, error) {
	var book models.Book
	err := r.DB.First(&book, id).Error
	return book, err
}

// GetBookWithLock 查询Book信息并获取行级锁
func (r *repository) GetBookWithLock(id uint) (models.Book, error) {
	var book models.Book
	err := r.DB.Raw(`
	select 
		id, title, price, stock 
	from books 
		where id =? and deleted_at is null 
	for update`, id).Scan(&book).Error
	return book, err
}

func (r *repository) ListBooks(f Filters) ([]*models.Book, error) {
	books := make([]*models.Book, 0, f.Limit())
	err := r.DB.Limit(f.Limit()).Offset(f.Offset()).Find(&books).Error
	return books, err
}

func (r *repository) DeductionStock(bookID uint, qty uint) error {
	var q models.Book
	err := r.DB.First(&q, bookID).Error
	if err != nil {
		return err
	}

	if q.Stock < qty {
		return ErrUnderstock
	}

	q.Stock -= qty

	return r.DB.Save(&q).Error
}

// TxUserBuyBook 购物图书事物
func (r *repository) TxUserBuyBook(userID uint, bookID uint) error {
	return r.execTx(func(tx Repository) error {
		book, err := tx.GetBookWithLock(bookID)
		if err != nil {
			return err
		}

		err = tx.CreateOrder(&models.Order{
			TotalAmount: book.Price,
			UserID:      userID,
			OrderDetails: []models.OrderDetail{
				{
					Qty:    1,
					Amount: book.Price,
					BookID: book.ID,
				},
			},
		})
		if err != nil {
			return err
		}

		err = tx.DeductionStock(bookID, 1)
		if err != nil {
			return err
		}

		return nil
	})
}
