package utils

import (
	"mime"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	ValidType       = regexp.MustCompile(`(image|video)\/[\s\S]*`)
	imageIndexRegex = regexp.MustCompile(`\(\d+\)`)
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

func GetImageIndex(filename string) (int64, error) {
	bracketedIndex := imageIndexRegex.FindString(filename)
	index := strings.Trim(bracketedIndex, "()")
	return strconv.ParseInt(index, 10, 32)
}

func FilesSortFunc(i, j int, v []string) bool {
	idx1, err := GetImageIndex(v[i])
	if err != nil {
		return false
	}
	idx2, err := GetImageIndex(v[j])
	if err != nil {
		return false
	}
	return idx1 < idx2
}
