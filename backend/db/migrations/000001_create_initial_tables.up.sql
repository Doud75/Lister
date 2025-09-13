CREATE TABLE bands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    band_id INT REFERENCES bands(id) ON DELETE CASCADE,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ... (CREATE TABLE bands et CREATE TABLE users restent avant) ...

-- Supprime l'ancienne définition de la table `songs` si elle existe dans ce fichier
-- CREATE TABLE songs (...);

-- Voici la nouvelle définition complète :
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    band_id INT NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    duration_seconds INT,
    tempo INT,
    song_key VARCHAR(10),
    lyrics TEXT,
    chords TEXT,
    album_name VARCHAR(255),
    instrumentation JSONB,
    notes TEXT,
    links JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE interludes (
    id SERIAL PRIMARY KEY,
    band_id INT REFERENCES bands(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    speaker VARCHAR(100),
    script TEXT,
    duration_seconds INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE setlists (
    id SERIAL PRIMARY KEY,
    band_id INT REFERENCES bands(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(7) DEFAULT '#FFFFFF',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE setlist_items (
    id SERIAL PRIMARY KEY,
    setlist_id INT NOT NULL REFERENCES setlists(id) ON DELETE CASCADE,
    position INT NOT NULL,
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('song', 'interlude')),
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    interlude_id INT REFERENCES interludes(id) ON DELETE CASCADE,
    notes TEXT,
    transition_duration_seconds INT DEFAULT 0,
    CONSTRAINT chk_item_is_defined CHECK (
        (item_type = 'song' AND song_id IS NOT NULL AND interlude_id IS NULL)
            OR
        (item_type = 'interlude' AND interlude_id IS NOT NULL AND song_id IS NULL)
    ),
   CONSTRAINT unique_position_in_setlist UNIQUE (setlist_id, position)
);

CREATE INDEX idx_songs_band_id_title ON songs(band_id, title);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_songs_updated_at
    BEFORE UPDATE ON songs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();