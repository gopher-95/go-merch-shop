CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    item_name VARCHAR(50) NOT NULL,
    quantity INT DEFAULT 1,
    UNIQUE(user_id, item_name)
);