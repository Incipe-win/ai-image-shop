package model

import (
	"gorm.io/gorm"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"    // 待支付
	OrderStatusPaid      OrderStatus = "paid"       // 已支付
	OrderStatusShipped   OrderStatus = "shipped"    // 已发货
	OrderStatusCompleted OrderStatus = "completed"  // 已完成
	OrderStatusCancelled OrderStatus = "cancelled"  // 已取消
)

type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	OrderSN     string         `gorm:"uniqueIndex;not null;size:64" json:"order_sn"`
	TotalAmount float64        `gorm:"not null;type:decimal(10,2)" json:"total_amount"`
	Status      OrderStatus    `gorm:"not null;size:20;default:'pending'" json:"status"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User       User        `gorm:"foreignKey:UserID" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrderID        uint           `gorm:"not null;index" json:"order_id"`
	ProductName    string         `gorm:"not null;size:100" json:"product_name"`
	ProductImageURL string         `gorm:"not null;size:500" json:"product_image_url"`
	DesignImageURL string         `gorm:"not null;size:500" json:"design_image_url"`
	Size           string         `gorm:"not null;size:10" json:"size"`
	Color          string         `gorm:"not null;size:30" json:"color"`
	Price          float64        `gorm:"not null;type:decimal(10,2)" json:"price"`
	Quantity       int            `gorm:"not null" json:"quantity"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Order Order `gorm:"foreignKey:OrderID" json:"-"`
}

func (OrderItem) TableName() string {
	return "order_items"
}