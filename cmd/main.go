package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tender-manager/config"
	"tender-manager/internal/app"
	"tender-manager/internal/app/bidsservice"
	"tender-manager/internal/app/employeeservice"
	"tender-manager/internal/app/tenderservice"

	gen "tender-manager/internal/generated"
	"tender-manager/internal/repository"
	"time"
)

func main() {
	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg, err := config.Parse()
	if err != nil {
		log.Error("Could not parse ENVs, ", err)
	}

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	router := chi.NewRouter()

	storage, err := repository.New(cfg)
	if err != nil {
		log.Error("Could not setup storage, ", err)
	}

	employeeCases := employeeservice.New(storage)
	tenderUseCases := tenderservice.New(storage, *employeeCases)
	bidUseCases := bidsservice.New(storage, *employeeCases, log)
	appStatus := app.New(&app.Status{})

	type CombinedUseCases struct {
		*tenderservice.Client
		*bidsservice.BidsClient
		app.StatusApp
	}

	combinedUseCases := CombinedUseCases{
		tenderUseCases,
		bidUseCases,
		appStatus,
	}
	h := http.Server{
		Addr:         cfg.ServerAddress,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Handler:      gen.HandlerFromMux(gen.NewStrictHandler(combinedUseCases, nil), router),
	}

	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Error("Listen And Serve", err)
		}
	}()

	log.Info("server running on", "address:", cfg.ServerAddress)
	<-ctx.Done()

	closeCtx, _ := context.WithTimeout(context.Background(), time.Second)
	if err := h.Shutdown(closeCtx); err != nil {
		log.Error("Http server close")
	}

}
