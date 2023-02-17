package domain

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"
)

const (
	OrderByName int = iota
	OrderByDate
)

type DirectoryEnt struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	Loved        bool      `json:"loved"`
	Thumbnail    string    `json:"thumbnail"`
	LastModified time.Time `json:"lastModified"`
}

type Directory struct {
	gorm.Model
	Path      string `gorm:"unique;not null"`
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

	FindAllRange(ctx context.Context, take, skip, order int) (*[]Directory, error)

	Update(ctx context.Context, path, name, thumbnail string) (Directory, error)
}

type ListingService interface {
	CountDirectories(ctx context.Context) (int64, error)

	ListAllDirectories(ctx context.Context) (*[]DirectoryEnt, error)

	ListAllDirectoriesLike(ctx context.Context, name string) (*[]DirectoryEnt, error)

	ListAllDirectoriesRange(ctx context.Context, take, skip, order int) (*[]DirectoryEnt, error)
}

type ListingHandler interface {
	ListAllDirectories() http.HandlerFunc
}
