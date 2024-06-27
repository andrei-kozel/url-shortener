package main

import (
	"os"

	"github.com/andrei-kozel/go-utils/utils/prettylog"
	"github.com/andrei-kozel/url-shortener/internal/config"
	"github.com/andrei-kozel/url-shortener/internal/storage/sqlite"
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

	log.Info("initializing storage...", "path", cfg.StoragePath)
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	_ = storage
}
