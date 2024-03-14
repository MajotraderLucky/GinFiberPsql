CREATE TABLE fio_data (
    id SERIAL PRIMARY KEY,
    name TEXT,
    surname TEXT,
    patronymic TEXT,
    age INTEGER,
    gender TEXT,
    nationality TEXT,
    error_reason TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX idx_fio_data_name ON fio_data (name);