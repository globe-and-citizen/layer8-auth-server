ALTER TABLE user_metadata
    ADD COLUMN last_email_verified_at TIMESTAMPTZ,
    ADD COLUMN last_phone_verified_at TIMESTAMPTZ;