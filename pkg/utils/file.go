package utils

import (
	"mime"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	THUMBNAIL_NAME = ".thumb.webp"
)

var (
	ValidType       = regexp.MustCompile(`(image|video)\/[\s\S]*`)
	ImageIndexRegex = regexp.MustCompile(`\(\d+\)`)
)

// ValidFile checks if file is eligible for viewing (isn't "hidden")
// or isn't the directory thumbnail
func ValidFile(filename string) bool {
	return filename != THUMBNAIL_NAME && !strings.HasPrefix(filename, ".")
}

func IsVideo(mime string) bool {
	return strings.HasPrefix(mime, "video")
}

func IsImage(mime string) bool {
	return strings.HasPrefix(mime, "image")
}

func IsImagePath(path string) bool {
	return strings.HasPrefix(mime.TypeByExtension(filepath.Ext(path)), "image")
}

func GetImageIndex(filename string) (int64, error) {
	bracketedIndex := ImageIndexRegex.FindString(filename)
	index := strings.Trim(bracketedIndex, "()")
	return strconv.ParseInt(index, 10, 32)
}
