CREATE TABLE IF NOT EXISTS artists (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS albums (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    artist_id  INTEGER NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    year       INTEGER,
    cover_path TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (artist_id, name)
);

CREATE TABLE IF NOT EXISTS tracks (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    album_id     INTEGER REFERENCES albums(id) ON DELETE SET NULL,
    artist_id    INTEGER REFERENCES artists(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    track_number INTEGER,
    disc_number  INTEGER,
    duration_sec INTEGER,
    file_path    TEXT NOT NULL UNIQUE,
    format       TEXT,
    bitrate      INTEGER,
    size_bytes   INTEGER,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS tracks_artist_id_idx ON tracks(artist_id);
CREATE INDEX IF NOT EXISTS tracks_album_id_idx  ON tracks(album_id);
CREATE INDEX IF NOT EXISTS albums_artist_id_idx ON albums(artist_id);
