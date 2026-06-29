package handler

import (
	"log/slog"
	"net/http"

	"github.com/GulzhanKarakul/subscription-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc        service.SubscriptionService
	log           *slog.Logger
}

func NewHandler(
	svc service.SubscriptionService,
	log *slog.Logger,
) *Handler {
	return &Handler{
		svc: svc,
		log:           log,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/subscriptions", func(r chi.Router) {
			// create
			r.Post("/", h.createSubscription)
			r.Get("/", h.getSubscriptionsList)
			r.Get("/total", h.calculateTotal)
			r.Get("/{id}", h.getSubscriptionByID)
			r.Put("/{id}", h.updateSubscription)
			r.Delete("/{id}", h.deleteSubscription)
		})
	})

	return r
}
