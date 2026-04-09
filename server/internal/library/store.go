package library

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// --- Artists ---

func (s *Store) UpsertArtist(ctx context.Context, name string) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO artists (name) VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, name).Scan(&id)
	return id, err
}

func (s *Store) ListArtists(ctx context.Context) ([]Artist, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name, created_at FROM artists ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []Artist
	for rows.Next() {
		var a Artist
		if err := rows.Scan(&a.ID, &a.Name, &a.CreatedAt); err != nil {
			return nil, err
		}
		artists = append(artists, a)
	}
	return artists, rows.Err()
}

func (s *Store) GetArtist(ctx context.Context, id int64) (*Artist, error) {
	var a Artist
	err := s.db.QueryRow(ctx, `
		SELECT id, name, created_at FROM artists WHERE id = $1
	`, id).Scan(&a.ID, &a.Name, &a.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

// --- Albums ---

func (s *Store) UpsertAlbum(ctx context.Context, artistID int64, name string, year *int) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO albums (artist_id, name, year)
		VALUES ($1, $2, $3)
		ON CONFLICT (artist_id, name) DO UPDATE SET year = COALESCE(EXCLUDED.year, albums.year)
		RETURNING id
	`, artistID, name, year).Scan(&id)
	return id, err
}

func (s *Store) ListAlbumsByArtist(ctx context.Context, artistID int64) ([]Album, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, artist_id, name, year, cover_path, created_at
		FROM albums WHERE artist_id = $1 ORDER BY year, name
	`, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		var a Album
		if err := rows.Scan(&a.ID, &a.ArtistID, &a.Name, &a.Year, &a.CoverPath, &a.CreatedAt); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, rows.Err()
}

func (s *Store) GetAlbum(ctx context.Context, id int64) (*Album, error) {
	var a Album
	err := s.db.QueryRow(ctx, `
		SELECT id, artist_id, name, year, cover_path, created_at
		FROM albums WHERE id = $1
	`, id).Scan(&a.ID, &a.ArtistID, &a.Name, &a.Year, &a.CoverPath, &a.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

// --- Tracks ---

func (s *Store) TrackExists(ctx context.Context, filePath string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM tracks WHERE file_path = $1)
	`, filePath).Scan(&exists)
	return exists, err
}

func (s *Store) InsertTrack(ctx context.Context, t Track) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO tracks
			(album_id, artist_id, title, track_number, disc_number, file_path, format, size_bytes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (file_path) DO NOTHING
	`, t.AlbumID, t.ArtistID, t.Title, t.TrackNumber, t.DiscNumber,
		t.FilePath, t.Format, t.SizeBytes)
	return err
}

func (s *Store) ListTracksByAlbum(ctx context.Context, albumID int64) ([]Track, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, album_id, artist_id, title, track_number, disc_number,
		       duration_sec, file_path, format, bitrate, size_bytes, created_at
		FROM tracks WHERE album_id = $1
		ORDER BY disc_number NULLS FIRST, track_number NULLS FIRST, title
	`, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTracks(rows)
}

func (s *Store) ListTracksByArtist(ctx context.Context, artistID int64) ([]Track, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, album_id, artist_id, title, track_number, disc_number,
		       duration_sec, file_path, format, bitrate, size_bytes, created_at
		FROM tracks WHERE artist_id = $1
		ORDER BY title
	`, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTracks(rows)
}

func (s *Store) GetTrack(ctx context.Context, id int64) (*Track, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, album_id, artist_id, title, track_number, disc_number,
		       duration_sec, file_path, format, bitrate, size_bytes, created_at
		FROM tracks WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tracks, err := scanTracks(rows)
	if err != nil || len(tracks) == 0 {
		return nil, err
	}
	return &tracks[0], nil
}

func (s *Store) SearchTracks(ctx context.Context, q string) ([]Track, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, album_id, artist_id, title, track_number, disc_number,
		       duration_sec, file_path, format, bitrate, size_bytes, created_at
		FROM tracks
		WHERE title ILIKE $1
		ORDER BY title
		LIMIT 50
	`, fmt.Sprintf("%%%s%%", q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTracks(rows)
}

// RandomTracks returns n random tracks (used by AutoDJ).
func (s *Store) RandomTracks(ctx context.Context, n int) ([]Track, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, album_id, artist_id, title, track_number, disc_number,
		       duration_sec, file_path, format, bitrate, size_bytes, created_at
		FROM tracks
		ORDER BY RANDOM()
		LIMIT $1
	`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTracks(rows)
}

func scanTracks(rows pgx.Rows) ([]Track, error) {
	var tracks []Track
	for rows.Next() {
		var t Track
		if err := rows.Scan(
			&t.ID, &t.AlbumID, &t.ArtistID, &t.Title,
			&t.TrackNumber, &t.DiscNumber, &t.DurationSec,
			&t.FilePath, &t.Format, &t.Bitrate, &t.SizeBytes, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, rows.Err()
}
