package listing

import (
	"fuu/v/pkg/common"
	"fuu/v/pkg/domain"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
)

type Handler struct {
	service domain.ListingService
}

// Handle the top level folder listing.
// Folders can be ordered by name or modified time.
func (h *Handler) ListAllDirectories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		dirs := new([]domain.DirectoryEnt)

		if filterBy != "" {
			dirs, err = h.service.ListAllDirectoriesLike(filterBy)
		} else {
			dirs, err = h.service.ListAllDirectoriesRange(pageSize, (page-1)*pageSize)
		}

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		count, err := h.service.CountDirectories()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		/*
			if r.URL.Query().Get("fetchBy") == "date" {
				sort.SliceStable(files, func(i, j int) bool {
					return files[i].CreatedAt.After(files[j].CreatedAt)
				})
			}
		*/

		paginator := common.Paginator[domain.DirectoryEnt]{
			Items:      dirs,
			PageSize:   pageSize,
			TotalItems: count,
		}

		body, err := json.Marshal(paginator.Get(page))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
