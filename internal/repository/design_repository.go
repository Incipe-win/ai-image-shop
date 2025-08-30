package repository

import (
	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"gorm.io/gorm"
)

type DesignRepository struct {
	db *gorm.DB
}

func NewDesignRepository(db *gorm.DB) *DesignRepository {
	return &DesignRepository{db: db}
}

func (r *DesignRepository) Create(design *model.Design) error {
	return r.db.Create(design).Error
}

func (r *DesignRepository) FindByUserID(userID uint) ([]model.Design, error) {
	var designs []model.Design
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&designs).Error
	return designs, err
}

func (r *DesignRepository) FindByUserIDAndCategory(userID uint, category string) ([]model.Design, error) {
	var designs []model.Design
	err := r.db.Where("user_id = ? AND category = ?", userID, category).Order("created_at DESC").Find(&designs).Error
	return designs, err
}

func (r *DesignRepository) FindByID(id uint) (*model.Design, error) {
	var design model.Design
	err := r.db.First(&design, id).Error
	return &design, err
}

func (r *DesignRepository) Delete(id uint) error {
	return r.db.Delete(&model.Design{}, id).Error
}
