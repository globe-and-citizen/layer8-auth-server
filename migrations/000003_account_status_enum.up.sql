-- Create the enum type
CREATE TYPE account_status AS ENUM (
    'owing',
    'zeroed',
    'overpaid'
    );

-- Remove the old default first
ALTER TABLE client_balance
    ALTER COLUMN status DROP DEFAULT;

-- Convert the existing column
ALTER TABLE client_balance
    ALTER COLUMN status
        TYPE account_status
        USING status::account_status;

-- Set the new default
ALTER TABLE client_balance
    ALTER COLUMN status
        SET DEFAULT 'zeroed'::account_status;