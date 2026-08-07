-- Create the enum type
CREATE TYPE account_status AS ENUM (
    'owing',
    'zeroed',
    'overpaid'
    );

-- Convert the existing column
ALTER TABLE client_balance
    ALTER COLUMN status
        TYPE account_status
        USING status::account_status;

-- Set the default
ALTER TABLE client_balance
    ALTER COLUMN status
        SET DEFAULT 'zeroed';
