-- Remove the default
ALTER TABLE client_balance
    ALTER COLUMN status
        DROP DEFAULT;

-- Convert back to TEXT
ALTER TABLE client_balance
    ALTER COLUMN status
        TYPE TEXT
        USING status::TEXT;

-- Drop the enum type
DROP TYPE account_status;