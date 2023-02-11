package domain

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type DirectoryEnt struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Loved     bool   `json:"loved"`
	Thumbnail string `json:"thumbnail"`
}

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

type ListingRepository interface {
	Count(ctx context.Context) (int64, error)

	Create(ctx context.Context, path, name, thumbnail string) (Directory, error)

	FindByName(ctx context.Context, name string) (Directory, error)

	FindAllByName(ctx context.Context, name string) (*[]Directory, error)

	Delete(ctx context.Context, path string) (Directory, error)

	FindAll(ctx context.Context) (*[]Directory, error)

	FindAllRange(ctx context.Context, take, skip int) (*[]Directory, error)

	Update(ctx context.Context, path, name, thumbnail string) (Directory, error)
}

type ListingService interface {
	CountDirectories() (int64, error)

	ListAllDirectories() (*[]DirectoryEnt, error)

	ListAllDirectoriesLike(name string) (*[]DirectoryEnt, error)

	ListAllDirectoriesRange(take, skip int) (*[]DirectoryEnt, error)
}

type ListingHandler interface {
	ListAllDirectories() http.HandlerFunc
}
