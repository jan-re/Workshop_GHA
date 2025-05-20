package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	targetServiceEnvName = "TARGET_SERVICE"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))
)

func main() {
	targetService, ok := os.LookupEnv(targetServiceEnvName)
	if !ok {
		logger.Error("Target service could not be retrieved from env. Shutting down.", slog.String("envName", targetServiceEnvName))
		os.Exit(1)
	}

	logger.Info("Application starting. Creating request emitter.")

	go handleProbes()

	go emit(targetService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	gotSignal := <-quit
	signal.Stop(quit)

	logger.Info("Termination signal received. Application shutting down.", slog.String("signal", gotSignal.String()))
}
