-- SQL schema for the worker-service jobs table.
-- This table is used for idempotency and tracking the status of blockchain operations.

CREATE TABLE IF NOT EXISTS jobs (
    id            TEXT PRIMARY KEY,
    event_id      TEXT UNIQUE NOT NULL, -- Ensures we never process the same event twice
    type          TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'pending',
    tx_hashes     TEXT NOT NULL DEFAULT '{}', -- PostgreSQL text array: {tx1,tx2}
    retry_count   INT  NOT NULL DEFAULT 0,
    error_message TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_jobs_event_id ON jobs (event_id);
CREATE INDEX IF NOT EXISTS idx_jobs_status   ON jobs (status);
