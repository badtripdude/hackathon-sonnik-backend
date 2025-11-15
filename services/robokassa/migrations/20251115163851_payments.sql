-- +goose Up
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- уникальный ID операции
    user_id UUID NOT NULL,                         -- ID пользователя (может не существовать в другой БД)
    amount NUMERIC(10,2) NOT NULL,                -- сумма платежа
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',   -- валюта
    status VARCHAR(20) NOT NULL DEFAULT 'created',-- статус платежа: created, pending, success, failed
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);

-- +goose Down
DROP TABLE payments;
