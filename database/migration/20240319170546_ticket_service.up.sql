CREATE TABLE IF NOT EXISTS tickets (
    id BIGINT PRIMARY KEY,
    capacity BIGINT,
    region VARCHAR(255),
    event_date TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ticket_details (
    id BIGINT PRIMARY KEY,
    ticket_id BIGINT,
    level VARCHAR(255),
    stock BIGINT,
    base_price DOUBLE PRECISION,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (ticket_id) REFERENCES tickets(ID)
);