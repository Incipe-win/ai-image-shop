package repository

import (
	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) AddItem(item *model.CartItem) error {
	// 检查是否已经存在相同的商品配置
	var existingItem model.CartItem
	err := r.db.Where("user_id = ? AND product_id = ? AND design_id = ? AND size = ? AND color = ?",
		item.UserID, item.ProductID, item.DesignID, item.Size, item.Color).First(&existingItem).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		return r.db.Create(item).Error
	} else if err != nil {
		return err
	}

	// 存在，更新数量
	existingItem.Quantity += item.Quantity
	return r.db.Save(&existingItem).Error
}

func (r *CartRepository) GetCartByUserID(userID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	err := r.db.Preload("Product").Preload("Design").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (r *CartRepository) UpdateQuantity(id uint, quantity int) error {
	return r.db.Model(&model.CartItem{}).Where("id = ?", id).Update("quantity", quantity).Error
}

func (r *CartRepository) RemoveItem(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.CartItem{}).Error
}

func (r *CartRepository) ClearCartByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.CartItem{}).Error
}

func (r *CartRepository) FindByID(id uint) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.Preload("Product").Preload("Design").First(&item, id).Error
	return &item, err
}