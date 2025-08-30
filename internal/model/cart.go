package model

import (
	"gorm.io/gorm"
	"time"
)

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	DesignID  uint           `gorm:"not null" json:"design_id"`
	Size      string         `gorm:"not null;size:10" json:"size"`
	Color     string         `gorm:"not null;size:30" json:"color"`
	Quantity  int            `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
	Design  Design  `gorm:"foreignKey:DesignID" json:"design"`
}

func (CartItem) TableName() string {
	return "cart_items"
}