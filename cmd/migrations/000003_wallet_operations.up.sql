CREATE TABLE IF NOT EXISTS wallet_operations (
    id bigserial PRIMARY KEY,
    wallet_id UUID NOT NULL,    
    balance NUMERIC(20, 2) NOT NULL DEFAULT 0,
    operation_type INT NOT NULL CHECK (operation_type IN (0, 1)),
    amount NUMERIC(20, 2) NOT NULL CHECK (amount > 0),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);