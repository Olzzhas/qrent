CREATE TABLE IF NOT EXISTS stations (
    id SERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_organizations
        FOREIGN KEY (org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE
);

CREATE TRIGGER stations_update_timestamp
BEFORE UPDATE ON stations
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

CREATE INDEX IF NOT EXISTS idx_stations_org_id
    ON stations(org_id);