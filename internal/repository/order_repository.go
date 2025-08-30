package repository

import (
	"fmt"
	"time"

	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) CreateWithItems(order *model.Order, items []model.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 为每个订单项设置 OrderID
		for i := range items {
			items[i].OrderID = order.ID
		}

		// 批量创建订单项
		if err := tx.CreateInBatches(items, 100).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *OrderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("OrderItems").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("OrderItems").First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) FindByOrderSN(orderSN string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("OrderItems").Where("order_sn = ?", orderSN).First(&order).Error
	return &order, err
}

func (r *OrderRepository) UpdateStatus(id uint, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) GenerateOrderSN() string {
	return fmt.Sprintf("ORD%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}