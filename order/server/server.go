package server

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/agumiroff/BigTechProject/order/v1/internal/handler/order"
	orderV1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func StartHTTPServer(h *order.OrderHandler, t time.Duration, port int) (*http.Server, error) {
	orderServer, err := orderV1.NewServer(h)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
		return nil, err
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", strconv.Itoa(port)),
		Handler:           r,
		ReadHeaderTimeout: t,
	}

	return server, nil
}
