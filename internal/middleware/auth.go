package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gopher-95/go-merch-shop/internal/service"
)

type AuthMiddleWare struct {
	jwt *service.JWT
}

func NewAuthMiddleware(jwt *service.JWT) *AuthMiddleWare {
	return &AuthMiddleWare{
		jwt: jwt,
	}
}

func (m *AuthMiddleWare) Handler(next http.Handler) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonResponseError(w, http.StatusUnauthorized, "неверный токен")
			return
		}

		token := ""

		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			jsonResponseError(w, http.StatusUnauthorized, "неправильный заголовок запроса")
			return
		}

		userID, err := m.jwt.Validate(token)
		if err != nil {
			jsonResponseError(w, http.StatusUnauthorized, "неправильный токен")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFunc)
}

// Функция возваращет error JSON зпрос
func jsonResponseError(w http.ResponseWriter, statuscode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": message,
	})
}
