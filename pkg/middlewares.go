package pkg

import (
	"fmt"
	"fuu/v/pkg/common"
	config "fuu/v/pkg/config"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	thumbsFormat = strings.ToLower(config.Instance().ImageOptimizationFormat)
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

// Middleware for allowing the serve of thumbnails as they're saved as file
// without extension. By rule thumbnails are AVIF pictures, so a Content-Type
// header is set.
func serveThumbnail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/"+thumbsFormat)
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
