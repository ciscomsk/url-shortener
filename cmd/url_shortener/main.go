package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// init logger: slog
	env := cfg.Env
	log := setupLogger(env)
	//log = log.With(slog.String("env", env))
	log.Info("Starting url-shortener", slog.String("env", env), slog.String("version", "123"))
	log.Debug("Debug messages are enabled")
	log.Error("Error message test")

	// init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	//id, err := storage.SaveURL("https://google.com", "google")
	//if err != nil {
	//	log.Error("failed to save url", sl.Err(err))
	//	os.Exit(1)
	//}
	//
	//log.Info("saved url", slog.Int64("id", id))
	//
	//id, err = storage.SaveURL("https://google.com", "google")
	//if err != nil {
	//	log.Error("failed to save url", sl.Err(err))
	//	os.Exit(1)
	//}
	//
	_ = storage

	// init router: chi + chi render
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	//router.Use(middleware.RealIP)
	//router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: run server

	// 1-11-00
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
		//log = slog.New(
		//	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		//)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
