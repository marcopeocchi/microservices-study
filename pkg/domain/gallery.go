package domain

import (
	"context"
	"net/http"
)

type Content struct {
	List          []string `json:"list"`
	Avif          []string `json:"avif"`
	AvifAvailable bool     `json:"avifAvailable"`
}

type DirectoryRepository interface {
	FindByPath(ctx context.Context, path string) (Content, error)
}

type DirectoryHandler interface {
	DirectoryContent() http.HandlerFunc
}
