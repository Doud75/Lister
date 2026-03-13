DROP INDEX IF EXISTS unique_default_band;
ALTER TABLE band_users DROP COLUMN IF EXISTS is_default;
