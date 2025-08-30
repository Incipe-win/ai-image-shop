package model

import (
	"gorm.io/gorm"
	"time"
)

type Design struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"size:200" json:"title"`                    // 创意名称
	Description string         `gorm:"type:text" json:"description"`             // 创意描述
	Prompt      string         `gorm:"type:text;not null" json:"prompt"`         // AI生成提示词
	ImageURL    string         `gorm:"not null;size:500" json:"image_url"`       // 创意图片URL
	Style       string         `gorm:"size:50" json:"style"`                     // 艺术风格
	Category    string         `gorm:"size:50;default:'general'" json:"category"` // 作品分类
	Tags        string         `gorm:"type:text" json:"tags"`                    // 标签，逗号分隔
	IsPublished bool           `gorm:"default:false" json:"is_published"`        // 是否已发布到商店
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联用户
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Design) TableName() string {
	return "designs"
}
