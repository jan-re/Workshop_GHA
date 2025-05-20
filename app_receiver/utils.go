package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func probeHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(getHelloWorldString()))
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"weather": "cloudy",
	}

	bs, err := json.Marshal(resp)
	if err != nil {
		panic(err.Error())
	}

	logger.Info("Request received. Sending positive response.", slog.String("userAgent", r.UserAgent()), slog.String("uri", r.RequestURI), slog.String("method", r.Method))

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

// getHelloWorldString only exists so that we can have a functional unit test for our pipeline.
func getHelloWorldString() string {
	return "Hello world"
}
