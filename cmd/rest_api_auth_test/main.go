package main

import (
	"net/http"
	"os"
	"restapiauthtest/internal/config"
	"restapiauthtest/internal/http_server/handlers/registration"
	"restapiauthtest/internal/storage/postgresql"
	"restapiauthtest/lib/logger"
	"restapiauthtest/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.MustLoad()
	log := logger.InitLogger(cfg.Env)

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
	// token, err := auth.Login(storage, "me", "password")
	// if err != nil {
	// 	log.Error("Could not login", sl.Err(err))
	// 	os.Exit(1)
	// }
	// log.Info("Logged in", sl.Token(token))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/registration", registration.Registration(log, storage))

	http.ListenAndServe(":8080", router)
	//пост запрос на регистрацию
	//
	// TODO: init handlers

	// TODO: /ping method

	defer func() {
		if err = storage.Close(); err != nil { //defered connection closure
			log.Error("Could not close storage", sl.Err(err))
			os.Exit(1)
		}
		log.Info("Storage closed")
	}()
}
