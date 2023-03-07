package gallery

import (
	"encoding/json"
	"fuu/v/internal/domain"
	"net/http"
)

type Handler struct {
	repo domain.DirectoryRepository
}

func (h *Handler) DirectoryContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer func() {
			r.Body.Close()
		}()

		dir := r.URL.Query().Get("dir")
		content, err := h.repo.FindByPath(r.Context(), dir)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(content)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}
}
