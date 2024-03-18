package main

import (
	"log/slog"
	"os"
	"restapiauthtest/internal/auth"
	"restapiauthtest/internal/config"
	"restapiauthtest/internal/storage/postgresql"
	"restapiauthtest/lib/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := initLogger(cfg.Env)
	storage, err := postgresql.New(cfg.StoragePath)
	defer func() {
		if err = storage.Close(); err != nil { //defered connection closure
			log.Error("Could not close storage", sl.Err(err))
			os.Exit(1)
		}
		log.Info("Storage closed")
	}()
	if err != nil {
		log.Error("Could not init storage", sl.Err(err))
		os.Exit(1)
	}
	// err = auth.RegisterNewUser(storage, "admin", "password")
	// if err != nil {
	// 	log.Error("Could not register new user", sl.Err(err))
	// 	os.Exit(1)
	// }
	_, err = auth.Login(storage, "admin", "password")
	if err != nil {
		log.Error("Could not login", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Logged in")

	// TODO: init logger
	// TODO: init storage
	// TODO: init server
	// TODO: init handlers
	// TODO: hash password & send to storage

	// TODO: create table with sessions, save unique token and lifetime

	// TODO: /ping method
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
