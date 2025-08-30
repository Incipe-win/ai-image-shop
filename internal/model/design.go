package model

import (
	"gorm.io/gorm"
	"time"
)

type Design struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Prompt    string         `gorm:"type:text;not null" json:"prompt"`
	ImageURL  string         `gorm:"not null;size:500" json:"image_url"`
	Style     string         `gorm:"size:50" json:"style"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联用户
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Design) TableName() string {
	return "designs"
}
