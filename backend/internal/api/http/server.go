package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-crypto-bot-clean/backend/internal/api/http/controllers"
	"go-crypto-bot-clean/backend/internal/api/http/middleware"
	"go-crypto-bot-clean/backend/internal/config"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	router    *mux.Router
	server    *http.Server
	container *config.Container
}

// NewServer creates a new HTTP server
func NewServer(container *config.Container) *Server {
	router := mux.NewRouter()
	
	// Create server
	server := &Server{
		router:    router,
		container: container,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", container.Config.App.Port),
			Handler: router,
		},
	}

	// Register middleware
	server.registerMiddleware()

	// Register routes
	server.registerRoutes()

	return server
}

// registerMiddleware registers middleware for the server
func (s *Server) registerMiddleware() {
	// Add common middleware
	s.router.Use(middleware.LoggingMiddleware)
	s.router.Use(middleware.RecoveryMiddleware)
	
	// Add CORS middleware
	s.router.Use(middleware.CORSMiddleware)

	// Add authentication middleware if enabled
	if s.container.Config.Auth.Enabled {
		s.router.Use(middleware.AuthMiddleware)
	}
}

// registerRoutes registers the API routes
func (s *Server) registerRoutes() {
	// Create controllers
	orderController := controllers.NewOrderController(s.container.OrderService)
	positionController := controllers.NewPositionController(s.container.PositionService)
	tradeController := controllers.NewTradeController(s.container.TradeService)

	// Register controller routes
	orderController.RegisterRoutes(s.router)
	positionController.RegisterRoutes(s.router)
	tradeController.RegisterRoutes(s.router)

	// Add health check endpoint
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server exited properly")
	return nil
}
