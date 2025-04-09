DROP TRIGGER IF EXISTS organizations_update_timestamp ON organizations;

DROP TABLE IF EXISTS organizations;

DROP FUNCTION IF EXISTS update_timestamp();