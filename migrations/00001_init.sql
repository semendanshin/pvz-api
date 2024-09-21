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

    received_at     TIMESTAMP    NOT NULL,
    storage_time    INTERVAL     NOT NULL,

    issued_at       TIMESTAMP,
    returned_at     TIMESTAMP,
    deleted_at      TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvz_orders;
-- +goose StatementEnd
