package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Структура хрнаит наш сервер
type Server struct {
	httpServer *http.Server
}

// Конструктор сервера
func NewServer(port string, router *chi.Mux) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
