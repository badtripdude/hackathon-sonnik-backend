package main

import (
	"context"
	repo "github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/repository/postgres"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/config"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/controller/http/server"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/service"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/auth"
)

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()
	db, err := postgres.NewPgClient(ctx, cfg)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	defer db.Close()

	privKey, err := auth.LoadPrivateKey(cfg.RSAPrivateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	pubKey, err := auth.LoadPublicKey(cfg.RSAPublicKeyPath)
	if err != nil {
		log.Fatalf("cannot load public key: %v", err)
	}

	queryTimeout := 5 * time.Second
	userRepo := repo.NewPgUserRepo(db, queryTimeout)
	refreshTokenRepo := repo.NewPgRefreshTokenRepo(db, queryTimeout)

	userSvc := service.NewUserService(*cfg, userRepo, refreshTokenRepo, privKey)

	r := server.NewRouter(userSvc, pubKey)

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
		log.Println("server starting on port", cfg.ServerPort)
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
