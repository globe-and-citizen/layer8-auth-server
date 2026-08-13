-- Remove the enum default
ALTER TABLE client_balance
    ALTER COLUMN status DROP DEFAULT;

-- Convert the enum back to text
ALTER TABLE client_balance
    ALTER COLUMN status
        TYPE VARCHAR
        USING status::text;

-- Restore the previous default
ALTER TABLE client_balance
    ALTER COLUMN status
        SET DEFAULT 'zeroed';

-- Remove the enum type
DROP TYPE account_status;