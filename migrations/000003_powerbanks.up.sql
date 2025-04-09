CREATE TABLE IF NOT EXISTS powerbanks (
    id SERIAL PRIMARY KEY,
    current_station_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_stations
        FOREIGN KEY (current_station_id)
        REFERENCES stations(id)
        ON DELETE CASCADE,
    CONSTRAINT chk_status
        CHECK (status IN ('rented', 'available', 'charging'))
);

CREATE TRIGGER powerbanks_update_timestamp
    BEFORE UPDATE ON powerbanks
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE INDEX IF NOT EXISTS idx_powerbanks_current_station_id
    ON powerbanks(current_station_id);