CREATE TYPE credit_type AS ENUM ('AUTO', 'MORTGAGE', 'COMMERCIAL');
CREATE TYPE credit_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TABLE IF NOT EXISTS credits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    bank_id UUID NOT NULL REFERENCES banks(id) ON DELETE CASCADE,
    min_payment DECIMAL(15, 2) NOT NULL,
    max_payment DECIMAL(15, 2) NOT NULL,
    term_months INTEGER NOT NULL,
    credit_type credit_type NOT NULL,
    status credit_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_credits_client_id ON credits(client_id);
CREATE INDEX idx_credits_status ON credits(status);