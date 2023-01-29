package models

import "gorm.io/gorm"

// Thumbnail DB Model
type Directory struct {
	gorm.Model
	Path      string `gorm:"primaryKey"`
	Name      string
	Loved     bool
	Thumbnail string
}
