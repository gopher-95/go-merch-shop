package models

// Структура описывает запрос пользователя на отправку монет
type SendCoinsRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
