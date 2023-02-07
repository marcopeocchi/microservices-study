package models

import (
	"time"

	"gorm.io/gorm"
)

// Thumbnail DB Model
type Directory struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Path      string         `gorm:"primaryKey;autoIncrement:false"`
	Name      string
	Loved     bool
	Thumbnail string
}
