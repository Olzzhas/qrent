DROP INDEX IF EXISTS idx_powerbanks_current_station_id;

DROP TRIGGER IF EXISTS powerbanks_update_timestamp ON powerbanks;

DROP TABLE IF EXISTS powerbanks;