-- Creating the fio_data table
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

-- Creating an index on the name field
CREATE INDEX idx_fio_data_name ON fio_data (name);

-- Creating a composite index on the name and surname fields
CREATE INDEX idx_fio_data_name_surname ON fio_data (name, surname);