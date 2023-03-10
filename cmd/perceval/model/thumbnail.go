package model

type Thumbnail struct {
	Thumbnail string
	Path      string
	Folder    string `gorm:"primaryKey"`
}
