package user

import (
	"context"
	"encoding/json"
	"fuu/v/pkg/common"
	"fuu/v/pkg/domain"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Handler struct {
	service domain.UserService
	logger  *zap.SugaredLogger
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())

		defer func() {
			cancel()
			r.Body.Close()
		}()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := loginRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.logger.Errorw("Decoding error", "error", err)
			return
		}

		token, err := h.service.Login(ctx, req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.logger.Warnw("Invalid credentials", "error", err)
			return
		}

		cookie := http.Cookie{
			Name:     common.TOKEN_COOKIE_NAME,
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(time.Minute * 30),
			Value:    *token,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)

		w.Write([]byte(*token))
	}
}

func (h *Handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:     common.TOKEN_COOKIE_NAME,
			HttpOnly: true,
			Expires:  time.UnixMilli(0),
			Value:    "",
		}
		http.SetCookie(w, &cookie)
	}
}
