package models

import (
 	"time"
 	"gorm.io/gorm"
)

type Item struct {
 	ID          uint       `gorm:"primaryKey" json:"id"`
 	CreatedAt   time.Time  `json:"createdAt"`
 	UpdatedAt   time.Time  `json:"updatedAt"`
 	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

 	Name        string     `gorm:"type:text;not null;index" json:"name"`
 	Description string     `gorm:"type:text" json:"description"`
 	Price       float64    `gorm:"type:real;index" json:"price"`
 	Status      string     `gorm:"type:text;default:active;index" json:"status"`
 	Version     uint       `gorm:"not null;default:1" json:"version"`
}
