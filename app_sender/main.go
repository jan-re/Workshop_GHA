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
	portEnvName = "PORT"
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

	logger.Info("Application starting. Creating server and router.", slog.String("port", port))

	mux := &http.ServeMux{}
	mux.Handle("GET /api/v1/weather", http.HandlerFunc(weatherHandler))
	mux.Handle("/live", http.HandlerFunc(probeHandler))
	mux.Handle("/ready", http.HandlerFunc(probeHandler))

	server := http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	logger.Info("Application has started successfully. Server is listening.", slog.String("port", port))

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("Server error. Application shutting down.", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	gotSignal := <-quit
	signal.Stop(quit)

	logger.Info("Termination signal received. Application shutting down.", slog.String("signal", gotSignal.String()))
}
