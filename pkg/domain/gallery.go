package domain

import (
	"context"
	"net/http"
)

type Content struct {
	Source        []string `json:"source"`
	Avif          []string `json:"avif"`
	WebP          []string `json:"webp"`
	AvifAvailable bool     `json:"avifAvailable"`
	WebPAvailable bool     `json:"webpAvailable"`
	Cached        bool     `json:"cached"`
}

type DirectoryRepository interface {
	FindByPath(ctx context.Context, path string) (Content, error)
}

type DirectoryHandler interface {
	DirectoryContent() http.HandlerFunc
}
