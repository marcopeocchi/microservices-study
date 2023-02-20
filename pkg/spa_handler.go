package pkg

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SpaHandler struct {
	entrypoint string
	filesystem *fs.FS
	routes     []string
}

func (s *SpaHandler) AddRoute(route string) *SpaHandler {
	s.routes = append(s.routes, route)
	return s
}

// Handler for serving a compiled react frontend: each client-side routes must be provided
func (s *SpaHandler) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		path := filepath.Clean(r.URL.Path)

		// basically all frontend routes are needed :/
		hasRoute := false
		for _, route := range s.routes {
			hasRoute = strings.HasPrefix(path, route)
			if hasRoute {
				break
			}
		}

		if path == "/" || hasRoute {
			path = s.entrypoint
		}

		path = strings.TrimPrefix(path, "/")

		file, err := (*s.filesystem).Open(path)

		if err != nil {
			if os.IsNotExist(err) {
				log.Println("file", path, "not found:", err)
				http.NotFound(w, r)
				return
			}
			log.Println("file", path, "cannot be read:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		w.Header().Set("Content-Type", contentType)

		if strings.HasPrefix(path, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=2592000")
		}

		stat, err := file.Stat()
		if err == nil && stat.Size() > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		}

		io.Copy(w, file)
	})
}