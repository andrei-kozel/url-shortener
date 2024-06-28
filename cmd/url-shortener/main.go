package main

import (
	"os"

	"github.com/andrei-kozel/go-utils/utils/prettylog"
	"github.com/andrei-kozel/url-shortener/internal/config"
	"github.com/andrei-kozel/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
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
}
