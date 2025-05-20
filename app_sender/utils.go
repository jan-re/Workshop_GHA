package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const targetEndpoint = "/api/v1/weather"

func emit(targetService string) {
	url := targetService + targetEndpoint

	counter := 1
	for {
		func() {
			defer func() {
				counter++
				time.Sleep(time.Second * 5)
			}()

			reqLogger := logger.With(slog.Int("requestNumber", counter))

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				panic("request is malformed: " + err.Error())
			}

			req.Header.Set("User-Agent", "Workshop_App_Sender")

			reqLogger.Info("Sending HTTP request.", slog.String("targetUrl", url))

			client := &http.Client{
				Timeout: time.Second * 10,
			}

			resp, err := client.Do(req)
			if err != nil {
				reqLogger.Error("Error received while sending request.", slog.String("err", err.Error()), slog.String("url", url))
				return
			}

			if resp.StatusCode != 200 {
				reqLogger.Error("Request sent. Received unexpected response status code.",
					slog.String("statusCode", resp.Status),
				)
				return
			}

			defer resp.Body.Close()

			var respBody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			if err != nil {
				reqLogger.Error("Request sent. Failed to decode response body.")
				return
			}

			reqLogger.Info("Request sent. Received positive response.",
				slog.String("statusCode", resp.Status),
				slog.String("weatherReport", respBody["weather"]),
			)
		}()
	}
}

func handleProbes() {
	mux := &http.ServeMux{}
	mux.Handle("/live", http.HandlerFunc(probeHandler))
	mux.Handle("/ready", http.HandlerFunc(probeHandler))

	server := http.Server{
		Addr:           ":54321",
		Handler:        mux,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			os.Exit(1)
		}
	}()
}

func probeHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(getHelloWorldString()))
}

// getHelloWorldString only exists so that we can have a functional unit test for our pipeline.
func getHelloWorldString() string {
	return "Hello world"
}
