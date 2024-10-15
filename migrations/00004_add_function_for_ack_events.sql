-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION get_pending_events(_limit INT)
    RETURNS SETOF events AS $$
DECLARE
    event_record events%ROWTYPE;
    count INT := 0;
BEGIN
    FOR event_record IN
        SELECT id, event_type, payload, created_at, sent_at
        FROM events
        WHERE sent_at IS NULL
        ORDER BY created_at
        LOOP
            IF pg_try_advisory_lock(event_record.id) THEN
                count := count + 1;
                RETURN NEXT event_record;

                IF count >= _limit THEN
                    EXIT;
                END IF;
            END IF;
        END LOOP;
    RETURN;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION mark_event_as_sent(event_id UUID)
    RETURNS VOID AS $$
BEGIN
    UPDATE events
    SET sent_at = NOW()
    WHERE id = event_id;

    PERFORM pg_advisory_unlock(event_id);
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose NO TRANSACTION
-- +goose Down
DROP FUNCTION IF EXISTS get_pending_events;
DROP FUNCTION IF EXISTS mark_event_as_sent;
