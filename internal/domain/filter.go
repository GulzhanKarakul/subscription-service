package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Limit       int
	Offset      int
}

type TotalFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        time.Time
	To          time.Time
}