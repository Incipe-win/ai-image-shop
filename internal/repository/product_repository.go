package repository

import (
	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("is_active = ?", true).Order("created_at DESC").Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) FindByCategory(category string) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("category = ? AND is_active = ?", category, true).Order("created_at DESC").Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}