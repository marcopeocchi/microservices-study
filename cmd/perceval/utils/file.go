package utils

import (
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	THUMBNAIL_NAME = ".thumb.webp"
)

var (
	ValidType = regexp.MustCompile(`(image|video)\/[\s\S]*`)
)

func ValidFile(filename string) bool {
	return !strings.HasPrefix(filename, ".")
}

func IsVideo(mime string) bool {
	return strings.HasPrefix(mime, "video")
}

func IsImage(mime string) bool {
	return strings.HasPrefix(mime, "image")
}

func IsImagePath(path string) bool {
	if !ValidFile(filepath.Base(path)) {
		return false
	}
	return IsImage(mime.TypeByExtension(filepath.Ext(path)))
}

func IsImageOrVideoPath(path string) bool {
	if !ValidFile(filepath.Base(path)) {
		return false
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	return IsImage(mimeType) || IsVideo(mimeType)
}
