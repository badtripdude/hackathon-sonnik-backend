// @title Robokassa Service API
// @version 1.0
// @description API для управления платежами через Robokassa
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/config"
	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/controller/http/server"
	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/repository/postgres"
	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/service"
	pg "github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/pkg/postgres"
)

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()
	db, err := pg.NewPgClient(ctx, cfg)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	defer db.Close()

	queryTimeout := 5 * time.Second
	paymentRepo := postgres.NewPgPaymentRepo(db, queryTimeout)
	paymentSvc := service.NewPaymentService(paymentRepo)

	r := server.NewRouter(paymentSvc)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           r,
		ReadTimeout:       cfg.ServerReadTimeout,
		WriteTimeout:      cfg.ServerWriteTimeout,
		IdleTimeout:       cfg.ServerIdleTimeout,
		ReadHeaderTimeout: cfg.ServerReadHeaderTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Robokassa service starting on port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed:", err)
		}
	}()

	<-stop
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("server shutdown failed:", err)
	} else {
		log.Println("server gracefully stopped")
	}
}
