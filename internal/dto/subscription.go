package dto

import (
	"time"

	"github.com/GulzhanKarakul/subscription-service/internal/domain"
)

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func ToSubscriptionResponse(s domain.Subscription) SubscriptionResponse {
	var endDate string
	if s.EndDate != nil {
		endDate = s.EndDate.Format("01-2006")
	}

	return SubscriptionResponse{
		ID:          s.ID.String(),
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID.String(),
		StartDate:   s.StartDate.Format("01-2006"),
		EndDate:     endDate,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}
}

func ToSubscriptionsResponse(subscriptions []domain.Subscription) []SubscriptionResponse {
	response := make([]SubscriptionResponse, 0, len(subscriptions))

	for _, subscription := range subscriptions {
		response = append(response, ToSubscriptionResponse(subscription))
	}

	return response
}

type TotalResponse struct {
	Total int `json:"total"`
}