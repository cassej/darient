-- Drop foreign keys
ALTER TABLE credits DROP CONSTRAINT IF EXISTS credits_client_id_fkey;
ALTER TABLE credits DROP CONSTRAINT IF EXISTS credits_bank_id_fkey;

-- Alter back to UUID
ALTER TABLE clients
    ALTER COLUMN id TYPE UUID USING gen_random_uuid(),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE banks
    ALTER COLUMN id TYPE UUID USING gen_random_uuid(),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE credits
    ALTER COLUMN id TYPE UUID USING gen_random_uuid(),
    ALTER COLUMN id SET DEFAULT gen_random_uuid(),
    ALTER COLUMN client_id TYPE UUID USING gen_random_uuid(),
    ALTER COLUMN bank_id TYPE UUID USING gen_random_uuid();

-- Recreate foreign keys
ALTER TABLE credits
    ADD CONSTRAINT credits_client_id_fkey
        FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE;

ALTER TABLE credits
    ADD CONSTRAINT credits_bank_id_fkey
        FOREIGN KEY (bank_id) REFERENCES banks(id) ON DELETE CASCADE;

-- Drop sequences
DROP SEQUENCE IF EXISTS clients_id_seq;
DROP SEQUENCE IF EXISTS banks_id_seq;
DROP SEQUENCE IF EXISTS credits_id_seq;