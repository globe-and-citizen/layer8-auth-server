ALTER TABLE client_balance
    DROP COLUMN balance_base_units,
    DROP COLUMN currency,
    ADD COLUMN balance_wei NUMERIC(78, 0) NOT NULL DEFAULT 0;

ALTER TABLE client_payment_receipt
    DROP COLUMN currency,
    DROP COLUMN sender_wallet_address,
    DROP COLUMN log_index;

DROP INDEX IF EXISTS uq_client_payment_receipt_tx_id_log_index;

ALTER TABLE client_payment_receipt
    RENAME COLUMN paid_amount TO paid_amount_wei;

UPDATE client_balance
SET balance_wei = 0,
    status = 'zeroed'::account_status;
