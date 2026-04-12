package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mattcarp12/mec/internal/auth"
	"github.com/mattcarp12/mec/internal/db"
	"github.com/mattcarp12/mec/internal/repository"
)

func main() {
	log.Println("Starting Identity Service...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admin:supersecretpassword@localhost:5432/identity?sslmode=disable"
	}

	migrate := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	if *migrate {
		if err := db.MigrateIdentityDB(dsn); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		os.Exit(0)
	}

	// Use a dummy 32-byte secret for local dev PASETO generation
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		tokenSecret = "12345678901234567890123456789012"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := db.NewPostgresDB(ctx, dsn)
	if err != nil {
		log.Fatalf("Fatal: Could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to Postgres!")

	// 1. Dependency Injection
	userRepo := repository.NewPostgresUserRepository(db)
	authService := auth.NewAuthService(userRepo)
	authHandler := auth.NewHandler(authService, tokenSecret)

	// 2. Go 1.22+ Standard Library Routing
	mux := http.NewServeMux()

	// Notice the HTTP methods included directly in the path string!
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	// 3. Server Configuration
	// Never use the default http.ListenAndServe() as it lacks timeouts and can easily
	// be targeted by Slowloris DDoS attacks. Always define a custom server.
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Listening on port 8080...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
