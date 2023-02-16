package pkg

import (
	"encoding/base64"

	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fuu/v/pkg/utils"

	"github.com/goccy/go-json"
	"github.com/marcopeocchi/fazzoletti/slices"
)

func listDirectoryContentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	dir := r.URL.Query().Get("dir")
	files, _ := os.ReadDir(filepath.Join(config.WorkingDir, dir))

	files = slices.Filter(files, func(file fs.DirEntry) bool {
		mimeType := mime.TypeByExtension(filepath.Ext(file.Name()))
		return utils.ValidType.MatchString(mimeType) && utils.ValidFile(file.Name())
	})

	res := make([]string, len(files))

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
