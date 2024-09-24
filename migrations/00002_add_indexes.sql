-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pvz_order_recipient_id ON pvz_orders (recipient_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pvz_order_pvz_id ON pvz_orders (pvz_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pvz_order_returned_at ON pvz_orders (returned_at);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pvz_order_received_at ON pvz_orders (received_at);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pvz_order_deleted_at_not_deleted ON pvz_orders (deleted_at) WHERE deleted_at IS NULL;

-- +goose NO TRANSACTION
-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_pvz_order_recipient_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_pvz_order_pvz_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_pvz_order_returned_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_pvz_order_received_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_pvz_order_deleted_at_not_deleted;
