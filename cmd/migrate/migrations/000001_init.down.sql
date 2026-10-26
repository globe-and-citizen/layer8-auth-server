-- Drop tables that reference others first
DROP TABLE IF EXISTS oauth_authorization_codes;

DROP TABLE IF EXISTS email_verification_data;
DROP TABLE IF EXISTS phone_number_verification_data;

DROP TABLE IF EXISTS client_payment_receipt;
DROP TABLE IF EXISTS client_balance;

DROP TABLE IF EXISTS user_metadata;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS clients;

DROP TABLE IF EXISTS zk_snarks_key_pairs;
