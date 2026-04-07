package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"github.com/gopher-95/go-merch-shop/internal/service"
)

// Хэндлер авторизации
type AuthHandler struct {
	service *service.AuthService
}

// Конструктор хэндлера авторизации
func NewAuthHanlder(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponseError(w, http.StatusMethodNotAllowed, "неправильный метод запроса")
		return
	}

	var authReq models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
		jsonResponseError(w, http.StatusBadRequest, "неправильное тело запроса")
		return
	}

	token, err := h.service.Login(r.Context(), authReq.Username, authReq.Password)
	if err != nil {
		jsonResponseError(w, http.StatusUnauthorized, "не получилось авторизоваться")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})

}

// Функция возваращет error JSON зпрос
func jsonResponseError(w http.ResponseWriter, statuscode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": message,
	})
}
