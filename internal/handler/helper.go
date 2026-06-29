package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GulzhanKarakul/subscription-service/internal/domain"
	"github.com/GulzhanKarakul/subscription-service/internal/dto"
	"github.com/google/uuid"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrSubscriptionNotFound):
		writeError(w, http.StatusNotFound, "not found")

	case errors.Is(err, domain.ErrSubscriptionAlreadyExist):
		writeError(w, http.StatusConflict, err.Error())

	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())

	default:
		h.log.Error("internal server error", "error", err.Error())
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func parseDate(s string) (time.Time, error) {
	date, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func subscriptionFromRequest(req dto.CreateSubscriptionRequest) (domain.Subscription, error) {
	userID, err := parseUUID(req.UserID)
	if err != nil {
			return domain.Subscription{}, fmt.Errorf("invalid user_id: %w", domain.ErrInvalidInput)
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
			return domain.Subscription{}, fmt.Errorf("invalid start_date: %w", domain.ErrInvalidInput)
	}

	var endDate *time.Time
	if req.EndDate != "" {
			t, err := parseDate(req.EndDate)
			if err != nil {
					return domain.Subscription{}, fmt.Errorf("invalid end_date: %w", domain.ErrInvalidInput)
			}
			endDate = &t
	}

	return domain.Subscription{
			ServiceName: req.ServiceName,
			Price:       req.Price,
			UserID:      userID,
			StartDate:   startDate,
			EndDate:     endDate,
	}, nil
}