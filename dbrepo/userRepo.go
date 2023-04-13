package dbrepo

import (
	"github.com/lightsaid/short-net/models"
)

func (r *repository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *repository) GetUserByID(id uint) (models.User, error) {
	var user models.User
	err := r.DB.Preload("Links").Limit(10).Offset(0).Find(&user, "id = ?", id).Error
	return user, err
}

func (r *repository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *repository) UpdateUserByID(id uint, user models.User) error {
	var q models.User
	err := r.DB.First(&q, id).Error
	if err != nil {
		return err
	}
	if user.Name != "" {
		q.Name = user.Name
	}
	if user.Avatar != "" {
		q.Avatar = user.Avatar
	}
	return r.DB.Save(&q).Error
}

func (r *repository) ActiveUserByID(id uint) error {
	var user models.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return err
	}
	return r.DB.Model(&user).Update("active", 1).Error
}
