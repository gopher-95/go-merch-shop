package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gopher-95/go-merch-shop/internal/service"
)

// Хэндлер описывает получение инфорамции о пользователе
type InfoHandler struct {
	service *service.InfoService
}

// Конструктор хэндлера информации
func NewInfoHandler(service *service.InfoService) *InfoHandler {
	return &InfoHandler{
		service: service,
	}
}

func (h *InfoHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponseError(w, http.StatusMethodNotAllowed, "неправильный метод запроса")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		jsonResponseError(w, http.StatusUnauthorized, "не авторизован")
		return
	}

	info, err := h.service.GetUserInfo(ctx, userID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			jsonResponseError(w, http.StatusGatewayTimeout, "таймаут запроса")
			return
		}
		jsonResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}
