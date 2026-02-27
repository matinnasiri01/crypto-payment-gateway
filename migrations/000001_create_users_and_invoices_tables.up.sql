CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    wallet VARCHAR(255) NOT NULL,
    balance FLOAT DEFAULT 0,
    currency VARCHAR(50) DEFAULT 'eth',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(255) PRIMARY KEY,
    creator_id VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    amount FLOAT NOT NULL,
    currency VARCHAR(50) DEFAULT 'eth',
    status VARCHAR(50) NOT NULL DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);


ALTER TABLE invoices 
ADD CONSTRAINT check_status_values 
CHECK (status IN ('Completed', 'Pending', 'Failed'));