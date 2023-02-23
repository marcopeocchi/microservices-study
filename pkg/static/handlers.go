package static

import (
	"encoding/base64"

	"mime"
	"net/http"
	"path/filepath"
	"strings"

	config "fuu/v/pkg/config"
)

func StreamVideoFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	path := r.URL.Query().Get("path")
	pathBytes, err := base64.URLEncoding.DecodeString(path)

	path = string(pathBytes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(path, config.Instance().WorkingDir) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))

	w.Header().Add("Content-Type", mimeType)
	http.ServeFile(w, r, path)
}
