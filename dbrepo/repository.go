package dbrepo

import (
	"fmt"

	"github.com/lightsaid/short-net/models"
	"gorm.io/gorm"
)

type Repository interface {
	// user repo
	CreateUser(user *models.User) error
	GetUserByID(id uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUserByID(id uint, user models.User) error

	// link repo
	CreateLink(link *models.Link) error
	GetLinkByID(id uint) (models.Link, error)
	GetLinkByHash(shortHash string) (models.Link, error)
	UpdateLinkByID(id uint, link models.Link) error
	DeleteLinkByID(id uint) error
	ListLinksByUserID(userID uint, f Filters) ([]*models.Link, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

// execTx 定义一个公共执行事务函数
func (r *repository) execTx(fn func(Repository) error) error {
	tx := r.DB.Begin()
	repo := NewRepository(tx)

	err := fn(repo)
	if err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("execTx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit().Error
}
