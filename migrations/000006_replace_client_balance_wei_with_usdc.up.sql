ALTER TABLE client_balance
    DROP COLUMN balance_wei,
    ADD COLUMN balance_base_units NUMERIC(78, 0) NOT NULL DEFAULT 0,
    ADD COLUMN currency TEXT NOT NULL DEFAULT 'usdc';

ALTER TABLE client_payment_receipt
    RENAME COLUMN paid_amount_wei TO paid_amount;

ALTER TABLE client_payment_receipt
    ADD COLUMN currency TEXT NOT NULL DEFAULT 'usdc',
    ADD COLUMN sender_wallet_address TEXT;

UPDATE client_payment_receipt
SET currency = 'pol';

UPDATE client_balance
SET balance_base_units = 0,
    currency = 'usdc',
    status = 'zeroed'::account_status;
