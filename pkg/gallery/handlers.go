package gallery

import (
	"context"
	"encoding/json"
	"fuu/v/pkg/domain"
	"net/http"
)

type Handler struct {
	repo domain.DirectoryRepository
}

func (h *Handler) DirectoryContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithCancel(context.Background())

		defer func() {
			r.Body.Close()
			cancel()
		}()

		dir := r.URL.Query().Get("dir")
		content, err := h.repo.FindByPath(ctx, dir)

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
