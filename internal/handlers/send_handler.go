package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"github.com/gopher-95/go-merch-shop/internal/service"
)

// Хэндлер обмена монетами
type SendHandler struct {
	service *service.SendCoinsService
}

// Конструктор хэнделра обмена монетами
func NewSendHandler(service *service.SendCoinsService) *SendHandler {
	return &SendHandler{
		service: service,
	}
}

func (h *SendHandler) SendCoins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponseError(w, http.StatusMethodNotAllowed, "неправильный запрос")
		return
	}

	fromUserID, ok := r.Context().Value("userID").(int)
	if !ok {
		jsonResponseError(w, http.StatusUnauthorized, "не удалось получить user_id пользователя")
		return
	}

	var req models.SendCoinsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponseError(w, http.StatusBadRequest, "неверный запрос")
		return
	}

	err := h.service.SendCoins(r.Context(), fromUserID, req.ToUser, req.Amount)
	if err != nil {
		jsonResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("монеты успешно отправлены"))

}
