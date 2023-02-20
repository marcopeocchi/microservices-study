package pkg

import (
	"fmt"
	"fuu/v/pkg/common"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Middleware for applying CORS policy for ALL hosts and for
// allowing ALL request headers.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

// Disable the file indexing of http.FileServer.
func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Middleware for serving a compiled react frontend: each client-side route
// must be provided
func embeddedSPAHandler(fs *fs.FS) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		path := filepath.Clean(r.URL.Path)

		// basically all frontend routes are needed :/
		if path == "/" ||
			strings.HasPrefix(path, "/gallery") ||
			strings.HasPrefix(path, "/login") {
			path = "index.html"
		}

		path = strings.TrimPrefix(path, "/")

		file, err := (*fs).Open(path)

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

// Middleware for allowing the serve of thumbnails as they're saved as file
// without extension. By rule thumbnails are AVIF pictures, so a Content-Type
// header is set.
func serveThumbnail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/avif")
		next.ServeHTTP(w, r)
	})
}

// Middleware for allowing only authenticated users to perform requests.
func authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if request came from localhost, if so disable security
		if os.Getenv("TESTING") != "" && strings.HasPrefix(r.RemoteAddr, "[::1]") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(common.TOKEN_COOKIE_NAME)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		if cookie == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWTSECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expiresAt, err := time.Parse(time.RFC3339, claims["expiresAt"].(string))

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if time.Now().After(expiresAt) {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			log.Println(err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
