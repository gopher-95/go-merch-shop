package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gopher-95/go-merch-shop/internal/service"
)

type BuyHandler struct {
	service *service.BuyService
}

func NewBuyHandler(service *service.BuyService) *BuyHandler {
	return &BuyHandler{
		service: service,
	}
}

func (h *BuyHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		jsonResponseError(w, http.StatusUnauthorized, "не авторизован")
		return
	}

	item := chi.URLParam(r, "item")
	if item == "" {
		jsonResponseError(w, http.StatusBadRequest, "товар не указан")
		return
	}

	err := h.service.BuyMerch(r.Context(), userID, item)
	if err != nil {
		msg := err.Error()
		if msg == "товар не найден" || msg == "недостаточно монет" {
			jsonResponseError(w, http.StatusBadRequest, msg)
			return
		}
		jsonResponseError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("мерч успешно приобретен"))
}
