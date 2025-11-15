package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/config"
)

type Server struct {
	cfg    *config.Config
	router http.Handler
	srv    *http.Server
}

func NewServer(cfg *config.Config, router http.Handler) *Server {
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
	}

	return &Server{
		cfg:    cfg,
		router: router,
		srv:    srv,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errc := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- err
		} else {
			errc <- nil
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Println("http server shutting down")
		return s.srv.Shutdown(shutdownCtx)
	case err := <-errc:
		if err != nil {
			log.Printf("http server error: %v\n", err)
		}
		return err
	}
}
