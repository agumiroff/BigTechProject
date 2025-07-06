package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"BigTechProject/hw1/cmd/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	httpPort          = "8080"
	urlParam          = "orderId"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	storage := models.NewOrderStorage()

	router := chi.NewRouter()

	// Middleware added
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// Routes
	router.Route("/api/order_service/v1/order", func(r chi.Router) {
		r.Get("/{orderId}", getOrderHandler(storage))
		r.Put("/{orderId}", updateOrderHandler(storage))
	})

	// Starting server
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           router,
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

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("there was a error accured when stopping server: %v\n", err)
	}

	log.Printf("Server successfully stopped")
}

func getOrderHandler(s *models.OrderStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getIdFromRequest(w, r)
		order := s.GetOrder(id)
		if order == nil {
			http.Error(w, fmt.Sprintf("Order for id %s not found", id), http.StatusNotFound)
			return
		}

		render.JSON(w, r, order)
	}
}

func updateOrderHandler(s *models.OrderStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getIdFromRequest(w, r)
		var orderUpdate models.Order
		if err := json.NewDecoder(r.Body).Decode(&orderUpdate); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		s.UpdateOrders(id, &orderUpdate)

		render.JSON(w, r, orderUpdate)
	}
}

func getIdFromRequest(w http.ResponseWriter, r *http.Request) string {
	id := chi.URLParam(r, urlParam)
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return ""
	}
	return id
}
