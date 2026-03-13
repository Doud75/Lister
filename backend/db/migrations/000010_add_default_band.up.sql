ALTER TABLE band_users ADD COLUMN is_default BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX unique_default_band ON band_users (user_id) WHERE is_default = TRUE;

UPDATE band_users bu
SET is_default = TRUE
FROM (
    SELECT DISTINCT ON (user_id) user_id, band_id
    FROM band_users
    ORDER BY user_id, band_id ASC
) first_band
WHERE bu.user_id = first_band.user_id AND bu.band_id = first_band.band_id;
