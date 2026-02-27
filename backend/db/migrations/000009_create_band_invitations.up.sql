CREATE TABLE band_invitations (
    id         SERIAL PRIMARY KEY,
    token      VARCHAR(255) NOT NULL UNIQUE,
    band_id    INT          NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    role       VARCHAR(50)  NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ  NOT NULL,
    max_uses   INT
);

CREATE INDEX idx_band_invitations_token ON band_invitations(token);
