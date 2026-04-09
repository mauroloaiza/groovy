package library

import "time"

type Artist struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Album struct {
	ID        int64     `json:"id"`
	ArtistID  int64     `json:"artist_id"`
	Name      string    `json:"name"`
	Year      *int      `json:"year,omitempty"`
	CoverPath *string   `json:"cover_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Track struct {
	ID          int64     `json:"id"`
	AlbumID     *int64    `json:"album_id,omitempty"`
	ArtistID    int64     `json:"artist_id"`
	Title       string    `json:"title"`
	TrackNumber *int      `json:"track_number,omitempty"`
	DiscNumber  *int      `json:"disc_number,omitempty"`
	DurationSec *int      `json:"duration_sec,omitempty"`
	FilePath    string    `json:"file_path"`
	Format      *string   `json:"format,omitempty"`
	Bitrate     *int      `json:"bitrate,omitempty"`
	SizeBytes   *int64    `json:"size_bytes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ScanResult is returned after a library scan.
type ScanResult struct {
	Scanned  int `json:"scanned"`
	Added    int `json:"added"`
	Skipped  int `json:"skipped"`
	Errors   int `json:"errors"`
}
