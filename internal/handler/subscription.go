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
// @Summary      Create subscription
// @Description  Create a new subscription for a user
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSubscriptionRequest true "Subscription data"
// @Success      201 {object} dto.SubscriptionResponse
// @Failure      400 {object} handler.ErrorResponse
// @Failure      500 {object} handler.ErrorResponse
// @Router       /subscriptions [post]
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
// @Summary      Get subscription by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id path string true "Subscription ID"
// @Success      200 {object} dto.SubscriptionResponse
// @Failure      400 {object} handler.ErrorResponse
// @Failure      404 {object} handler.ErrorResponse
// @Failure      500 {object} handler.ErrorResponse
// @Router       /subscriptions/{id} [get]
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
// @Summary      Update subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "Subscription ID"
// @Param        request body dto.UpdateSubscriptionRequest true "Updated data"
// @Success      200 {object} dto.SubscriptionResponse
// @Failure      400 {object} handler.ErrorResponse
// @Failure      404 {object} handler.ErrorResponse
// @Failure      500 {object} handler.ErrorResponse
// @Router       /subscriptions/{id} [put]
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

// @Summary      Delete subscription
// @Tags         subscriptions
// @Param        id path string true "Subscription ID"
// @Success      204
// @Failure      400 {object} handler.ErrorResponse
// @Failure      404 {object} handler.ErrorResponse
// @Failure      500 {object} handler.ErrorResponse
// @Router       /subscriptions/{id} [delete]
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

// @Summary      List subscriptions
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query string false "Filter by user ID"
// @Param        service_name query string false "Filter by service name"
// @Param        limit        query int    false "Limit (default 20)"
// @Param        offset       query int    false "Offset (default 0)"
// @Success      200 {array}  dto.SubscriptionResponse
// @Failure      500 {object} handler.ErrorResponse
// @Router       /subscriptions [get]
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

// @Summary      Calculate total cost
// @Tags         subscriptions
// @Produce      json
// @Param        from         query string false "From date MM-YYYY"
// @Param        to           query string false "To date MM-YYYY"
// @Param        user_id      query string false "Filter by user ID"
// @Param        service_name query string false "Filter by service name"
// @Success      200 {object} dto.TotalResponse
// @Failure      400 {object} handler.ErrorResponse
// @Failure      500 {object} handler.ErrorResponse
// @Router       /subscriptions/total [get]
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
