package main

import (
	"log/slog"
	"net/http"
	"os"
	"restapiauthtest/internal/auth"
	"restapiauthtest/internal/config"
	"restapiauthtest/internal/storage/postgresql"
	"restapiauthtest/lib/logger/sl"

	"github.com/go-chi/chi/v5"
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
	if err != nil {
		log.Error("Could not init storage", sl.Err(err))
		os.Exit(1)
	}

	// err = auth.RegisterNewUser(storage, "me", "password")
	// if err != nil {
	// 	log.Error("Could not register new user", sl.Err(err))
	// 	os.Exit(1)
	// }

	token, err := auth.Login(storage, "me", "password")
	if err != nil {
		log.Error("Could not login", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Logged in", sl.Token(token))

	r := chi.NewRouter()

	// Определяем обработчик для корневого URL
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет, мир!"))
	})

	http.ListenAndServe(":8080", r)

	// TODO: init server
	// TODO: init handlers

	// TODO: create table with sessions, save unique token and lifetime

	// TODO: /ping method

	defer func() {
		if err = storage.Close(); err != nil { //defered connection closure
			log.Error("Could not close storage", sl.Err(err))
			os.Exit(1)
		}
		log.Info("Storage closed")
	}()
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
