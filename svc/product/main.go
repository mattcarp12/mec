package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mattcarp12/mec/internal/auth"
	"github.com/mattcarp12/mec/internal/config"
	"github.com/mattcarp12/mec/internal/db"
)

func main() {
	log.Println("Starting Product Service...")

	// Initialize dependencies
	config.Init()
	db.InitDB()
	db.InitRedis()

	jwks := auth.NewJWKSClient(config.Get().JWTIssuer)

	err := jwks.Refresh(context.TODO())
	if err != nil {
		log.Fatalf("Failed to refresh JWKS: %v", err)
	}

	authMiddleware := auth.Middleware(jwks, config.Get().JWTIssuer, config.Get().JWTAudience)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/hello", Hello_Handler)
	mux.Handle("GET /api/v1/hello_p", authMiddleware(http.HandlerFunc(Hello_Handler)))

	// Server Configuration
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

func Hello_Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello from Product Service!"))
}
