package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	portEnvName    = "PORT"
	targetEndpoint = "/api/v1/weather"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))
)

func main() {
	port, ok := os.LookupEnv(portEnvName)
	if !ok {
		logger.Error("Port could not be retrieved from env. Shutting down.", slog.String("envName", portEnvName))
		os.Exit(1)
	}

	logger.Info("Application starting. Creating request emitter.")

	url := "http://localhost:" + port + targetEndpoint

	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error received while sending request.", slog.String("err", err.Error()), slog.String("url", url))
	}

	logger.Info("Request sent. Received response.", slog.String("statusCode", resp.Status))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	gotSignal := <-quit
	signal.Stop(quit)

	logger.Info("Termination signal received. Application shutting down.", slog.String("signal", gotSignal.String()))
}

func emit(port string) {
	url := "http://localhost:" + port + targetEndpoint

	for {
		func() {
			defer time.Sleep(time.Second * 5)

			// TODO Log sending
			resp, err := http.Get(url)
			if err != nil {
				logger.Error("Error received while sending request.", slog.String("err", err.Error()), slog.String("url", url))
				return
			}

			logger.Info("Request sent. Received response.", slog.String("statusCode", resp.Status))
			// TODO Log negatives
		}()
	}
}
