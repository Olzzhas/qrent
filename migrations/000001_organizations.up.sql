CREATE TABLE IF NOT EXISTS organizations (
     id SERIAL PRIMARY KEY,
     name text NOT NULL,
     location text NOT NULL,
     created_at timestamptz NOT NULL DEFAULT NOW()
);