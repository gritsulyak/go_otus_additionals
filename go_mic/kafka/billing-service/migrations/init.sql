CREATE TABLE IF NOT EXISTS inbox (
    message_id  TEXT PRIMARY KEY,
    received_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS billing (
    id         TEXT PRIMARY KEY,
    amount     BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
