ALTER TABLE users ADD COLUMN band_id INT REFERENCES bands(id) ON DELETE CASCADE;

UPDATE users u
SET band_id = (
    SELECT bu.band_id
    FROM band_users bu
    WHERE bu.user_id = u.id
    LIMIT 1
    );

DROP TABLE band_users;