CREATE TABLE IF NOT EXISTS stations (
    id SERIAL PRIMARY KEY,
    org_id integer NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_organizations
    FOREIGN KEY (org_id)
    REFERENCES organizations(id)
    ON DELETE CASCADE
);