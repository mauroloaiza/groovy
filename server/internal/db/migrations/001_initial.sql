CREATE TABLE IF NOT EXISTS artists (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS albums (
    id         BIGSERIAL PRIMARY KEY,
    artist_id  BIGINT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    year       INT,
    cover_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (artist_id, name)
);

CREATE TABLE IF NOT EXISTS tracks (
    id           BIGSERIAL PRIMARY KEY,
    album_id     BIGINT REFERENCES albums(id) ON DELETE SET NULL,
    artist_id    BIGINT REFERENCES artists(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    track_number INT,
    disc_number  INT,
    duration_sec INT,
    file_path    TEXT NOT NULL UNIQUE,
    format       TEXT,
    bitrate      INT,
    size_bytes   BIGINT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS tracks_artist_id_idx ON tracks(artist_id);
CREATE INDEX IF NOT EXISTS tracks_album_id_idx  ON tracks(album_id);
CREATE INDEX IF NOT EXISTS albums_artist_id_idx ON albums(artist_id);
