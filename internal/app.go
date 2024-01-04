package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/vadimpk/url-pinger/config"
	httpcontroller "github.com/vadimpk/url-pinger/internal/controller/http"
	"github.com/vadimpk/url-pinger/pkg/httpserver"
	logging "github.com/vadimpk/url-pinger/pkg/logger"
)

func Run(config *config.Config) {
	logger := logging.New(config.Log.Level)

	logger.Info("Starting application", "config", config)

	controller := httpcontroller.New(httpcontroller.Options{
		Logger: logger,
		Config: config,
	})

	server := httpserver.New(
		controller,
		httpserver.Port(config.HTTP.Port),
		httpserver.ReadTimeout(config.HTTP.ReadTimeout),
		httpserver.WriteTimeout(config.HTTP.WriteTimeout),
		httpserver.ShutdownTimeout(config.HTTP.ShutdownTimeout),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())

	case err := <-server.Notify():
		logger.Error("app - Run - server.Notify", "err", err)
	}

	err := server.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown", "err", err)
	}

	// drain all the other connections
}
