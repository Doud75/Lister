CREATE TABLE band_users (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    band_id INT NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    PRIMARY KEY (user_id, band_id)
);

INSERT INTO band_users (user_id, band_id, role)
SELECT id, band_id, 'admin' FROM users WHERE band_id IS NOT NULL;

ALTER TABLE users DROP COLUMN band_id;