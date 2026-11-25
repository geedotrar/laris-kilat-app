CREATE TABLE IF NOT EXISTS flash_items (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    price NUMERIC(12, 2) NOT NULL,
    stock INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stock_reservations (
    id SERIAL PRIMARY KEY,
    item_id VARCHAR(50) REFERENCES flash_items(id),
    user_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'RESERVED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO flash_items (id, title, price, stock)
VALUES ('item-sneaker-01', 'LarisKilat Midnight Runner X', 1499000, 50)
ON CONFLICT (id) DO UPDATE SET stock = 50;
