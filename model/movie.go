package model

import "time"

type Movie struct {
	ID          string    `gorm:"id,primaryKey" json:"id"`
	Title       string    `gorm:"title" json:"title"`
	Description string    `gorm:"description" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
}
