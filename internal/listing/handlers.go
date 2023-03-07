package listing

import (
	"fuu/v/internal/domain"
	"fuu/v/pkg/common"
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

		var orderBy int
		fetchBy := r.URL.Query().Get("fetchBy")

		if fetchBy == "date" {
			orderBy = domain.OrderByDate
		} else {
			orderBy = domain.OrderByName
		}

		filterBy := r.URL.Query().Get("filter")

		dirs := new([]domain.DirectoryEnt)

		if filterBy != "" {
			dirs, err = h.service.ListAllDirectoriesLike(r.Context(), filterBy)
		} else {
			dirs, err = h.service.ListAllDirectoriesRange(r.Context(), pageSize, (page-1)*pageSize, orderBy)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count, err := h.service.CountDirectories(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		paginator := common.Paginator[domain.DirectoryEnt]{
			Items:      dirs,
			PageSize:   pageSize,
			TotalItems: count,
		}

		body, err := json.Marshal(paginator.Get(page))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
