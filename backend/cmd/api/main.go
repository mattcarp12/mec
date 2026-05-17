package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mattcarp12/mec/internal/platform/config"
	"github.com/mattcarp12/mec/internal/platform/db"
	"github.com/mattcarp12/mec/internal/platform/logger"
	"github.com/mattcarp12/mec/internal/platform/middleware"
)

func main() {
	// =========================================================================
	// Initial Application Setup
	// =========================================================================

	// --- Load Configuration ---
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fatal config error: %v", err)
	}

	// --- Initialize Logger ---
	log := logger.New(cfg.Environment)
	log.Info().Str("env", cfg.Environment).Msg("Starting application")

	// --- Connect to Database ---
	dbConn, _, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Database initialization failed")
	}
	defer dbConn.Close()
	log.Println("Successfully connected to PostgreSQL")

	// =========================================================================
	// HTTP Router Setup
	// =========================================================================
	r := chi.NewRouter()

	// Essential SaaS Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger) // In Prod, swap this for a custom zerolog middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Health Check Route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

	// Protected API Routes (v1)
	r.Route("/api/v1", func(r chi.Router) {
		// Mount the PocketBase token introspection middleware
		// Note: In production, "http://localhost:8090" should come from your config struct!
		r.Use(middleware.PocketBaseAuth(cfg.PocketbaseURL))

		// A test endpoint to verify the context injection
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			// Extract the ID that our middleware injected into the request context
			userID := r.Context().Value(middleware.UserIDKey).(string)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "authenticated", "pocketbase_user_id": "` + userID + `"}`))
		})
	})

	// TODO: Mount domain routes here (e.g., /api/v1/products, /api/v1/orders)

	// =========================================================================
	// Server Start & Graceful Shutdown
	// =========================================================================
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Msg("API server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down API server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}
	log.Info().Msg("Server exited properly")
}
