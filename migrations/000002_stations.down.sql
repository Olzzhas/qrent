DROP INDEX IF EXISTS idx_stations_org_id;

DROP TRIGGER IF EXISTS stations_update_timestamp ON stations;

DROP TABLE IF EXISTS stations;