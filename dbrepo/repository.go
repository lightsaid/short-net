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
	ActiveUserByID(id uint) error

	// transaction

	// NOTE: 这里尝试一个新做法：通过 callback 函数解耦注册和发送邮件（严格意义来讲，也不是解耦）
	// 注册事物，user 是用户对象，cb 函数告知事物已经执行创建用户这一步，但并没有真正提及事物,
	// 等待发送邮件成功后再提交事物，因此可以在 callback 函数里根据错误信息提前做出响应
	// TxRegister(user *models.User, cb func(err error)) error
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
		return err
	}

	return tx.Commit().Error
}

// TxRegister 注册事物
// func (r *repository) TxRegister(user *models.User, callback func(err error)) error {
// 	return r.execTx(func(r Repository) error {
// 		err := r.CreateUser(user)
// 		callback(err)
// 		if err != nil {
// 			return err
// 		}

// 		 TODO:发送邮件
// 		return errors.New("故意而为之～")
// 	})
// }
