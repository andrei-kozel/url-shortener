package main

import (
	"net/http"
	"os"

	"github.com/andrei-kozel/go-utils/utils/prettylog"
	"github.com/andrei-kozel/url-shortener/internal/config"
	"github.com/andrei-kozel/url-shortener/internal/http-server/handlers/redirect"
	"github.com/andrei-kozel/url-shortener/internal/http-server/handlers/url/save"
	"github.com/andrei-kozel/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := prettylog.SetupLoggger(cfg.Env)

	log.Info("starting url-shortener...", "env", cfg.Env)
	log.Debug("debug messages are enabled...")

	// Here we initialize the storage
	log.Info("initializing storage...", "path", cfg.StoragePath)
	storage, err := sqlite.New(cfg.StoragePath)
	_ = storage
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	// Inititalize the router
	log.Info("initializing router...")
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// routes
	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting http server...", "port", cfg.HTTPServer.Address)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start http server", err)
		os.Exit(1)
	}

	log.Info("shutting down...")
}
