package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	router := chi.NewRouter()

	// Middleware added
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// Routes

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
