CREATE TABLE IF NOT EXISTS transactions(
    id SERIAL PRIMARY KEY,
    from_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    amount INT NOT NULL CHECK (amount >0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_from_user ON transactions(from_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_user ON transactions(to_user_id);