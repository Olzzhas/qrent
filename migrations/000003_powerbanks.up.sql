CREATE TABLE IF NOT EXISTS powerbanks (
    id SERIAL PRIMARY KEY,
    current_station_id integer NOT NULL,
    status text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_stations
    FOREIGN KEY (current_station_id)
    REFERENCES stations(id)
    ON DELETE CASCADE,
    CONSTRAINT chk_status
    CHECK (status IN ('rented', 'available', 'charging'))
);