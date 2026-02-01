-- Drop existing foreign keys
ALTER TABLE credits DROP CONSTRAINT IF EXISTS credits_client_id_fkey;
ALTER TABLE credits DROP CONSTRAINT IF EXISTS credits_bank_id_fkey;

-- Create sequences
CREATE SEQUENCE IF NOT EXISTS clients_id_seq;
CREATE SEQUENCE IF NOT EXISTS banks_id_seq;
CREATE SEQUENCE IF NOT EXISTS credits_id_seq;

-- Alter clients table
ALTER TABLE clients
    ALTER COLUMN id TYPE BIGINT USING 1,
    ALTER COLUMN id SET DEFAULT nextval('clients_id_seq');

-- Alter banks table
ALTER TABLE banks
    ALTER COLUMN id TYPE BIGINT USING 1,
    ALTER COLUMN id SET DEFAULT nextval('banks_id_seq');

-- Alter credits table
ALTER TABLE credits
    ALTER COLUMN id TYPE BIGINT USING 1,
    ALTER COLUMN id SET DEFAULT nextval('credits_id_seq'),
    ALTER COLUMN client_id TYPE BIGINT USING 1,
    ALTER COLUMN bank_id TYPE BIGINT USING 1;

-- Recreate foreign keys
ALTER TABLE credits
    ADD CONSTRAINT credits_client_id_fkey
        FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE RESTRICT;

ALTER TABLE credits
    ADD CONSTRAINT credits_bank_id_fkey
        FOREIGN KEY (bank_id) REFERENCES banks(id) ON DELETE RESTRICT;

-- Reset sequences to start from 1
SELECT setval('clients_id_seq', 1, false);
SELECT setval('banks_id_seq', 1, false);
SELECT setval('credits_id_seq', 1, false);