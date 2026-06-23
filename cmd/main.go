package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	logger.Info("server started", "addr", ":8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("serve failed", "err", err)
		os.Exit(1)
	}
}