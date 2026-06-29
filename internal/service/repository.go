package service

import (
	"context"
	"time"

	"github.com/GulzhanKarakul/subscription-service/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *domain.Subscription) (domain.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	Update(ctx context.Context, subscription *domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]domain.Subscription, error)
	CalculateTotal(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int64, error)
}
