package dbrepo

import (
	"github.com/lightsaid/short-net/models"
)

func (r *repository) CreateLink(link *models.Link) error {
	return r.DB.Create(link).Error
}

func (r *repository) GetLinkByID(id uint) (models.Link, error) {
	var link models.Link
	err := r.DB.First(&link, id).Error
	return link, err
}

func (r *repository) GetLinkByHash(shortHash string) (models.Link, error) {
	var link models.Link
	err := r.DB.Where("short_hash = ?", shortHash).First(&link).Error
	return link, err
}

func (r *repository) UpdateLinkByID(id uint, link models.Link) error {
	var q models.Link
	err := r.DB.First(&q, id).Error
	if err != nil {
		return err
	}
	if link.LongURL != "" {
		q.LongURL = link.LongURL
	}
	return r.DB.Save(&q).Error
}
func (r *repository) DeleteLinkByID(id uint) error {
	return r.DB.Delete(&models.Link{ID: id}).Error
}

func (r *repository) ListLinksByUserID(userID uint, f Filters) ([]*models.Link, error) {
	links := make([]*models.Link, 0, f.Limit())
	err := r.DB.Limit(f.Limit()).Offset(f.Offset()).Find(&links, "user_id = ?", userID).Error
	return links, err
}
