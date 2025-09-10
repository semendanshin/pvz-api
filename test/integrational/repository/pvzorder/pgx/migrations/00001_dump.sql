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

INSERT INTO pvz_orders (order_id, pvz_id, recipient_id, cost, weight, packaging, additional_film, received_at,
                        storage_time, issued_at, returned_at, deleted_at)
VALUES ('1', '1', '1', 1000, 1000, 'box', false, '2021-01-01T00:00:00Z', '0 years 0 mons 1 days 0 hours 0 mins 0.0 secs', NULL, NULL, NULL),
       ('2', '1', '1', 1000, 1000, 'bag', false, '2021-01-02T00:00:00Z', '0 years 0 mons 1 days 0 hours 0 mins 0.0 secs', NULL, NULL, NULL),

       ('3', '2', '1', 1000, 1000, 'film', false, '2021-01-03T00:00:00Z', '0 years 0 mons 1 days 0 hours 0 mins 0.0 secs', NULL, NULL, NULL),
       ('4', '1', '2', 1000, 1000, 'box', false, '2021-01-04T00:00:00Z', '0 years 0 mons 1 days 0 hours 0 mins 0.0 secs', NULL, NULL, NULL),
       ('5', '1', '2', 1000, 1000, 'box', false, '2021-01-01T00:00:00Z', '0 years 0 mons 1 days 0 hours 0 mins 0.0 secs', '2021-01-02T00:00:00Z',
        '2021-01-03T00:00:00Z', NULL),
       ('6', '1', '2', 1000, 1000, 'box', false, '2021-01-01T00:00:00Z', '0 years 0 mons 1 days 0 hours 0 mins 0.0 secs', '2021-01-02T00:00:00Z',
        '2021-01-03T00:00:00Z', NULL),
       ('7', '1', '2', 1000, 1000, 'box', false, '2021-01-01T00:00:00Z', '0 years 0 mons 1 days 0 hours 0 mins 0.0 secs', '2021-01-02T00:00:00Z',
        '2021-01-03T00:00:00Z', NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvz_orders;
-- +goose StatementEnd