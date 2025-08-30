package model

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	BasePrice   float64        `gorm:"not null;type:decimal(10,2)" json:"base_price"`
	Category    string         `gorm:"size:50" json:"category"`
	Brand       string         `gorm:"size:50" json:"brand"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	
	// 创意作品相关字段
	DesignID     *uint   `gorm:"index" json:"design_id"`           // 关联的设计ID
	CreatorID    *uint   `gorm:"index" json:"creator_id"`          // 创作者ID
	CreatorName  string  `gorm:"size:100" json:"creator_name"`     // 创作者名称
	DesignPrompt string  `gorm:"type:text" json:"design_prompt"`   // 设计提示词
	DesignStyle  string  `gorm:"size:50" json:"design_style"`      // 设计风格
	ImageURL     string  `gorm:"size:500" json:"image_url"`        // 作品图片URL
	
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Design  *Design `gorm:"foreignKey:DesignID" json:"design,omitempty"`
	Creator *User   `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

func (Product) TableName() string {
	return "products"
}