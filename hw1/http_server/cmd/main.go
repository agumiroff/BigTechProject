package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	ordersV1 "github.com/agumiroff/BigTechProject/hw1/http_server/pkg/openapi/orders/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*ordersV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*ordersV1.Order),
	}
}

func (s *OrderStorage) GetOrder(orderId string) *ordersV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderId]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) UpdateOrder(orderId string, order *ordersV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[orderId] = order
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) GetOrderByID(_ context.Context, params ordersV1.GetOrderByIDParams) (ordersV1.GetOrderByIDRes, error) {
	order := h.storage.GetOrder(params.OrderID)
	if order == nil {
		return &ordersV1.NotFoundError{
			Code:    404,
			Message: "Order for orderId'" + params.OrderID + "'not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) UpdateOrderByID(_ context.Context, req *ordersV1.Order, params ordersV1.UpdateOrderByIDParams) (ordersV1.UpdateOrderByIDRes, error) {
	order := &ordersV1.Order{
		OrderUUID:       params.OrderID,
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      req.TotalPrice,
		TransactionUUID: req.TransactionUUID,
		PaymentMethod:   req.PaymentMethod,
		Status:          req.Status,
	}

	h.storage.UpdateOrder(order.OrderUUID, order)

	return order, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *ordersV1.GenericErrorStatusCode {
	return &ordersV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: ordersV1.GenericError{
			Code:    ordersV1.NewOptInt(http.StatusInternalServerError),
			Message: ordersV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	orderServer, err := ordersV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()
	// Middleware added
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	// Starting server
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	go func() {
		log.Printf("Starting server port %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server starting error %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("shutting down server")

	// Context with timeout to stop server
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("there was a error accured when stopping server: %v\n", err)
	}

	log.Printf("Server successfully stopped")
}
