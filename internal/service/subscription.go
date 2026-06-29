package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/GulzhanKarakul/subscription-service/internal/domain"
	"github.com/google/uuid"
)

type subscriptionService struct {
	repo SubscriptionRepository
	log  *slog.Logger
}

func NewSubscriptionService(repo SubscriptionRepository, log *slog.Logger) SubscriptionService {
	return &subscriptionService{repo: repo, log: log}
}

func validateSubscription(subscription *domain.Subscription) error {
	if subscription.ServiceName == "" {
		return errors.New("service name is required")
	}
	if subscription.Price <= 0 {
		return errors.New("price must be positive number")
	}
	if subscription.UserID == uuid.Nil {
		return errors.New("user id is nil")
	}
	if subscription.StartDate.IsZero() {
		return errors.New("start date is required")
	}
	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		return errors.New("end date is lower then start date")
	}
	return nil
}

// Create subscriptionService method
func (s *subscriptionService) Create(ctx context.Context, subscription *domain.Subscription) (domain.Subscription, error) {
	if err := validateSubscription(subscription); err != nil {
		return domain.Subscription{}, fmt.Errorf("subscriptionService.Create: %w", err)
	}

	sub, err := s.repo.Create(ctx, subscription)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("subscriptionService.Create: %w", err)
	}

	return sub, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("subscriptionService.GetByID: %w", err)
	}

	return sub, nil
}

// Update subscription
func (s *subscriptionService) Update(ctx context.Context, subscription *domain.Subscription) (domain.Subscription, error) {
	if err := validateSubscription(subscription); err != nil {
		return domain.Subscription{}, fmt.Errorf("subscriptionService.Update: %w", err)
	}

	sub, err := s.repo.Update(ctx, subscription)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("subscriptionService.Update: %w", err)
	}

	return sub, nil
}

// Delete subscription
func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("subscriptionService.Delete: %w", err)
	}
	return nil
}

// List subs
func (s *subscriptionService) List(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]domain.Subscription, error) {
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	subSl, err := s.repo.List(ctx, userID, serviceName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("subscriptionService.List: %w", err)
	}

	return subSl, nil
}

// CalculateTotal
func (s *subscriptionService) CalculateTotal(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int64, error) {
	if from.IsZero() || to.IsZero() {
		return 0, fmt.Errorf("subscriptionService.CalculateTotal: %w: from and to are required", domain.ErrInvalidInput)
	}
	if to.Before(from) {
		return 0, fmt.Errorf("subscriptionService.CalculateTotal: %w: from must be before to", domain.ErrInvalidInput)
	}

	total, err := s.repo.CalculateTotal(ctx, from, to, userID, serviceName)
	if err != nil {
		return 0, fmt.Errorf("subscriptionService.CalculateTotal: %w", err)
	}

	return total, nil
}
