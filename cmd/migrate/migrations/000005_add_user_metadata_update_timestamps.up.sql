ALTER TABLE user_metadata
    ADD COLUMN display_name_updated_at TIMESTAMPTZ,
    ADD COLUMN color_updated_at TIMESTAMPTZ,
    ADD COLUMN bio_updated_at TIMESTAMPTZ,
    ADD COLUMN email_verified_at TIMESTAMPTZ,
    ADD COLUMN phone_number_verified_at TIMESTAMPTZ;