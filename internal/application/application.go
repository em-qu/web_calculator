package application

import (
	"context"
   "log/slog"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"

	"github.com/em-qu/web_calculator/internal/config"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

// Функция запуска приложения
func (a *Application) Run() error {
	cfg := config.MustLoad()

   log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("starting web_calculator")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      nil, //router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")
	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server")
		return err
	}
	log.Info("server stopped")
   return nil
}

