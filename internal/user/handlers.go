package user

import (
	"encoding/json"
	"fuu/v/internal/domain"
	"fuu/v/pkg/common"
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

type singupRequest = loginRequest

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
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

		token, err := h.service.Login(r.Context(), req.Username, req.Password)
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

func (h *Handler) SingUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r.Body.Close()
		}()

		req := singupRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.logger.Errorw("Decoding error", "error", err)
			return
		}

		user, err := h.service.Create(
			r.Context(),
			req.Username,
			req.Password,
			domain.Standard,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.logger.Errorw("Bad request", "error", err, "req", req)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}
