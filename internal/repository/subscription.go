package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GulzhanKarakul/subscription-service/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *subscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Create subscription
func (r *subscriptionRepository) Create(
	ctx context.Context,
	subscription *domain.Subscription,
) (domain.Subscription, error) {
	const query = `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var s domain.Subscription
	err := r.db.QueryRowContext(
		ctx, query, subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate,
	).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.Subscription{}, domain.ErrSubscriptionAlreadyExist
		}
		return domain.Subscription{}, fmt.Errorf("subscriptionRepository.Crete: %w", err)
	}

	return s, nil
}

// GetByID
func (r *subscriptionRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (domain.Subscription, error) {
	const query = `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var s domain.Subscription
	err := r.db.QueryRowContext(
		ctx, query, id,
	).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Subscription{}, domain.ErrSubscriptionNotFound
		}
		return domain.Subscription{}, fmt.Errorf("subscriptionRepository.GetByID: %w", err)
	}

	return s, nil
}

// Update
func (r *subscriptionRepository) Update(
	ctx context.Context,
	subscription *domain.Subscription,
) (domain.Subscription, error) {
	const query = `
		UPDATE subscriptions
		SET service_name = $1,
			price = $2,
			start_date = $3,
			end_date = $4,
			updated_at = NOW()
		WHERE id = $5
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var s domain.Subscription
	err := r.db.QueryRowContext(
		ctx, query, subscription.ServiceName, subscription.Price, subscription.StartDate, subscription.EndDate, subscription.ID,
	).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Subscription{}, domain.ErrSubscriptionNotFound
		}
		return domain.Subscription{}, fmt.Errorf("subscriptionRepository.Update: %w", err)
	}

	return s, nil
}

// Delete
func (r *subscriptionRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	const query = `
		DELETE FROM subscriptions
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("subscriptionRepository.Delete: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

// List
func (r *subscriptionRepository) List(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	limit, offset int,
) ([]domain.Subscription, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)
	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, *userID)
		argIdx++
	}
	if serviceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argIdx))
		args = append(args, *serviceName)
		argIdx++
	}

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
	`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("subscriptionRepository.List: query subscriptions: %w", err)
	}
	defer rows.Close()

	result := make([]domain.Subscription, 0)
	for rows.Next() {
		var s domain.Subscription
		if err = rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("subscriptionRepository.List: get subscriptions list by user id: %w", err)
		}
		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("subscriptionRepository.List: scan subscriptions: %w", err)
	}

	return result, nil
}

// CalculateTotal
func (r *subscriptionRepository) CalculateTotal(
	ctx context.Context,
	from, to time.Time,
	userID *uuid.UUID,
	serviceName *string,
) (int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	conditions = append(conditions,
		fmt.Sprintf("start_date <= $%d", argIdx),
	)
	args = append(args, to)
	argIdx++

	conditions = append(conditions,
		fmt.Sprintf("(end_date IS NULL OR end_date >= $%d)", argIdx),
	)
	args = append(args, from)
	argIdx++

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, *userID)
		argIdx++
	}
	if serviceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argIdx))
		args = append(args, *serviceName)
		argIdx++
	}

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE ` + strings.Join(conditions, " AND ")

	var total int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrSubscriptionNotFound
		}
		return 0, fmt.Errorf("subscriptionRepository: get Subscriptions by id: %w", err)
	}

	return total, nil
}
