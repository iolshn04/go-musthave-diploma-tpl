package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/auth"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/response"
)

type AuthHandler struct {
	service   *auth.Service
	secretKey string
}

func NewAuthHandler(service *auth.Service, secretKey string) *AuthHandler {
	return &AuthHandler{
		service:   service,
		secretKey: secretKey,
	}
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondJSON(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Login == "" || req.Password == "" {
		response.RespondJSON(w, http.StatusBadRequest, "login and password are required")
		return
	}

	user, err := h.service.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		response.RespondJSON(w, http.StatusConflict, "login already exists")
		return
	}

	value := user.ID + "|" + middleware.Sign(user.ID, h.secretKey)

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	})

	response.RespondJSON(w, http.StatusOK, "User successfully registered")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondJSON(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Login == "" || req.Password == "" {
		response.RespondJSON(w, http.StatusBadRequest, "login and password are required")
		return
	}

	user, err := h.service.Authenticate(r.Context(), req.Login, req.Password)
	if err != nil {
		response.RespondJSON(w, http.StatusUnauthorized, "invalid login or password")
		return
	}

	value := user.ID + "|" + middleware.Sign(user.ID, h.secretKey)

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	})

	response.RespondJSON(w, http.StatusOK, "User successfully logged in")
}
