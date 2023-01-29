package pkg

import (
	"encoding/base64"
	"strconv"
	"time"

	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fuu/v/pkg/models"
	"fuu/v/pkg/utils"

	"github.com/goccy/go-json"
	"github.com/marcopeocchi/fazzoletti/slices"
)

// Handle the top level folder listing.
// Folders can be ordered by name or modified time.
func listDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	pageParam := r.URL.Query().Get("page")
	pageSizeParam := r.URL.Query().Get("pageSize")
	page := 1
	pageSize := 50

	var err error

	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			page = 1
		}
	}
	if pageSizeParam != "" {
		pageSize, err = strconv.Atoi(pageSizeParam)
		if err != nil {
			pageSize = 50
		}
	}

	filterBy := r.URL.Query().Get("filter")

	files := []models.Directory{}

	if filterBy != "" {
		db.Where("name LIKE ?", "%"+filterBy+"%").Find(&files)
	} else {
		db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&files)
	}

	if r.URL.Query().Get("fetchBy") == "date" {
		sort.SliceStable(files, func(i, j int) bool {
			return files[j].CreatedAt.After(files[i].CreatedAt)
		})
	}

	var count int64
	db.Table("directories").Count(&count)

	if err != nil {
		log.Fatal(err)
	}

	list := make([]Preview, len(files))

	for i, file := range files {
		list[i].Id = file.ID
		list[i].Name = file.Name
		list[i].Loved = file.Loved
		list[i].Thumbnail = file.Thumbnail
	}

	paginator := Paginator[Preview]{
		Items:      &list,
		TotalItems: count,
		PageSize:   pageSize,
	}

	res := paginator.Get(page)

	body, err := json.Marshal(res)

	if err != nil {
		log.Panicln(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

	// GC
	res = nil
	list = nil
	body = nil
}

func listDirectoryContentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	dir := r.URL.Query().Get("dir")
	files, _ := os.ReadDir(filepath.Join(config.WorkingDir, dir))

	files = slices.Filter(files, func(file fs.DirEntry) bool {
		mimeType := mime.TypeByExtension(filepath.Ext(file.Name()))
		return utils.ValidType.MatchString(mimeType) && utils.ValidFile(file.Name())
	})

	res := make([]DirectortyList, len(files))

	for i, file := range files {
		if !file.IsDir() {
			res[i] = file.Name()
		}
	}

	sort.SliceStable(res, func(i, j int) bool {
		idx1, err := utils.GetImageIndex(res[i])
		if err != nil {
			return false
		}
		idx2, err := utils.GetImageIndex(res[j])
		if err != nil {
			return false
		}
		return idx1 < idx2
	})

	body, err := json.Marshal(Response[DirectortyList]{
		List: res,
	})

	if err != nil {
		log.Fatalln(err)
	}

	w.Write(body)

	// GC
	res = nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	config := appContext.Value("config").(Config)

	req := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, _ := GenerateToken2String(
		[]byte(req.Password),
		[]byte(config.ServerSecret),
	)

	if req.Password == config.Masterpass {
		cookie := http.Cookie{
			Name:     "fuutoken",
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			Value:    token,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func streamVideoFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	path := r.URL.Query().Get("path")
	pathBytes, err := base64.URLEncoding.DecodeString(path)

	path = string(pathBytes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(path, config.WorkingDir) {
		log.Println(config.WorkingDir)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))

	w.Header().Add("Content-Type", mimeType)
	http.ServeFile(w, r, path)
}
