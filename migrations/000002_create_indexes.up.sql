CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);

CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);

CREATE INDEX idx_subscriptions_user_start_date ON subscriptions(user_id, start_date);