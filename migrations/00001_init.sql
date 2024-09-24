-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS pvz_orders
(
    order_id        VARCHAR(255) PRIMARY KEY,
    pvz_id          VARCHAR(255) NOT NULL,
    recipient_id    VARCHAR(255) NOT NULL,

    cost            INT          NOT NULL,
    weight          INT          NOT NULL,

    packaging       VARCHAR(255) NOT NULL,
    additional_film BOOLEAN      NOT NULL DEFAULT FALSE,

    received_at     timestamptz  NOT NULL,
    storage_time    INTERVAL     NOT NULL,

    issued_at       timestamptz,
    returned_at     timestamptz,
    deleted_at      timestamptz
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvz_orders;
-- +goose StatementEnd
