ALTER TABLE client_balance
    ADD CONSTRAINT fk_client_balance_client
        FOREIGN KEY (client_id)
            REFERENCES clients(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;

ALTER TABLE user_metadata
    ADD CONSTRAINT fk_user_metadata_user
        FOREIGN KEY (id)
            REFERENCES users (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;

ALTER TABLE client_payment_receipt
    ADD CONSTRAINT fk_client_payment_receipt_client
        FOREIGN KEY (client_id)
            REFERENCES clients (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;

ALTER TABLE email_verification_data
    ADD CONSTRAINT fk_email_verification_data_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;

ALTER TABLE phone_number_verification_data
    ADD CONSTRAINT fk_phone_number_verification_data_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;

ALTER TABLE oauth_authorization_codes
    ADD CONSTRAINT fk_oauth_authorization_codes_client
        FOREIGN KEY (client_id)
            REFERENCES clients(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;

ALTER TABLE oauth_authorization_codes
    ADD CONSTRAINT fk_oauth_authorization_codes_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;
