CREATE TABLE IF NOT EXISTS payments (
    id          TEXT PRIMARY KEY,
    amount      BIGINT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS outbox (
    id           BIGSERIAL PRIMARY KEY,
    aggregate_id TEXT NOT NULL,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_outbox_unprocessed
    ON outbox (id) WHERE processed_at IS NULL;
