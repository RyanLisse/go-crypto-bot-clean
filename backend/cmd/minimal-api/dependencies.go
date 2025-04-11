package main

import (
	// Standard library
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// Chi router and middleware
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	// Configuration
	"github.com/spf13/viper"

	// Logging
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// Database
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Validation
	"github.com/go-playground/validator/v10"

	// UUID generation
	"github.com/google/uuid"

	// Project imports
	"go-crypto-bot-clean/backend/internal/config"
)

// This file is used to import and initialize key dependencies
// It ensures that all required dependencies are included in the build
// and can be gradually incorporated into the minimal API

// initDependencies initializes key dependencies for demonstration purposes
// This function is not called in the minimal API but serves as a reference
func initDependencies() {
	// This function is not meant to be called
	// It's here to ensure dependencies are included in the build
	fmt.Println("Initializing dependencies (this is just a reference function)")

	// Initialize Chi router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize Logrus
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetLevel(logrus.InfoLevel)

	// Initialize Zap
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLogger, _ := zapConfig.Build()
	defer zapLogger.Sync()

	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Initialize SQLite with sqlx
	sqlxDB, _ := sqlx.Connect("sqlite3", ":memory:")
	defer sqlxDB.Close()

	// Initialize SQLite with standard library
	sqlDB, _ := sql.Open("sqlite3", ":memory:")
	defer sqlDB.Close()

	// Initialize GORM
	gormDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	// Initialize validator
	validate := validator.New()
	_ = validate.Struct(struct {
		Name string `validate:"required"`
	}{Name: "test"})

	// Generate UUID
	id := uuid.New()
	fmt.Println(id.String())

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use context in a request
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	_ = req
}
