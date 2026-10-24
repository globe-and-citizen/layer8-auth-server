ALTER TABLE client_balance
    ADD COLUMN unbilled_usage_bytes NUMERIC(78, 0) NOT NULL DEFAULT 0;
