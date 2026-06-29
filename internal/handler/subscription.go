package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GulzhanKarakul/subscription-service/internal/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// POST /api/v1/subscriptions - Create
func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	subscription, err := subscriptionFromRequest(req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	sub, err := h.svc.Create(r.Context(), &subscription)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToSubscriptionResponse(sub))
}

// GET /api/v1/subscriptions/{id} - GetByID
func (h *Handler) getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subID, err := parseUUID(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	sub, err := h.svc.GetByID(r.Context(), subID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToSubscriptionResponse(sub))
}

// PUT /api/v1/subscriptions/{id} - Update
func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subID, err := parseUUID(id) // 
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	subscription, err := subscriptionFromRequest(req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	subscription.ID = subID

	sub, err := h.svc.Update(r.Context(), &subscription)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToSubscriptionResponse(sub))
}

// DELETE /api/v1/subscriptions/{id} - Delete
func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subID, err := parseUUID(id) // 
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.svc.Delete(r.Context(), subID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/subscriptions - List
func (h *Handler) getSubscriptionsList(w http.ResponseWriter, r *http.Request) {
	var userID *uuid.UUID
	if id := r.URL.Query().Get("user_id"); id != "" {
		parsed, err := parseUUID(id)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &parsed
	}
	var serviceName *string
	if sn := r.URL.Query().Get("service_name"); sn != "" {
		serviceName = &sn
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	subscriptions, err := h.svc.List(r.Context(), userID, serviceName, limit, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToSubscriptionsResponse(subscriptions))
}

// GET /api/v1/subscriptions/total - CalculateTotal
func (h *Handler) calculateTotal(w http.ResponseWriter, r *http.Request) {
	var userID *uuid.UUID
	if id := r.URL.Query().Get("user_id"); id != "" {
		parsed, err := parseUUID(id)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &parsed
	}
	var serviceName *string
	if sn := r.URL.Query().Get("service_name"); sn != "" {
		serviceName = &sn
	}

	from := r.URL.Query().Get("from")
	if from == "" {
		writeError(w, http.StatusBadRequest, "from is required")
		return
	}
	fromDate, err := parseDate(from)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid from input")
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		writeError(w, http.StatusBadRequest, "to is required")
		return
	}
	toDate, err := parseDate(to)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid to input")
		return
	}

	total, err := h.svc.CalculateTotal(r.Context(), fromDate, toDate, userID, serviceName)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.TotalResponse{Total: total})
}
