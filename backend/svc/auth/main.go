package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mattcarp12/mec/internal/db"
)

func main() {
	log.Println("Starting Auth Service...")
	db.Init()

	mux := http.NewServeMux()

	// Notice the HTTP methods included directly in the path string!
	mux.HandleFunc("POST /api/v1/auth/register", Register_Handler)
	mux.HandleFunc("POST /api/v1/auth/login", Login_Handler)
	mux.HandleFunc("POST /api/v1/auth/refresh", Refresh_Handler)
	mux.HandleFunc("POST /api/v1/auth/logout", Logout_Handler)

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
