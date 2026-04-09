package library

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

var audioExts = map[string]bool{
	".mp3":  true,
	".flac": true,
	".ogg":  true,
	".m4a":  true,
	".aac":  true,
	".wav":  true,
	".opus": true,
	".wma":  true,
}

// Scanner walks a directory and upserts tracks into the store.
type Scanner struct {
	store *Store
}

func NewScanner(store *Store) *Scanner {
	return &Scanner{store: store}
}

func (s *Scanner) Scan(ctx context.Context, root string) ScanResult {
	var result ScanResult

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("scan walk error at %s: %v", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !audioExts[ext] {
			return nil
		}

		result.Scanned++

		// Skip if already indexed.
		exists, err := s.store.TrackExists(ctx, path)
		if err != nil {
			log.Printf("scan check exists %s: %v", path, err)
			result.Errors++
			return nil
		}
		if exists {
			result.Skipped++
			return nil
		}

		if err := s.indexFile(ctx, path, ext); err != nil {
			log.Printf("scan index %s: %v", path, err)
			result.Errors++
			return nil
		}

		result.Added++
		return nil
	})

	if err != nil {
		log.Printf("scan error: %v", err)
	}

	log.Printf("scan complete — scanned:%d added:%d skipped:%d errors:%d",
		result.Scanned, result.Added, result.Skipped, result.Errors)

	return result
}

func (s *Scanner) indexFile(ctx context.Context, path, ext string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	size := info.Size()

	// Read tags; fall back to filename if tags are missing.
	m, err := tag.ReadFrom(f)

	var (
		artistName  string
		albumName   string
		title       string
		year        *int
		trackNum    *int
		discNum     *int
	)

	if err != nil || m == nil {
		// No tags — use filename as title, "Unknown" for artist/album.
		title = strings.TrimSuffix(filepath.Base(path), ext)
		artistName = "Unknown Artist"
		albumName = "Unknown Album"
	} else {
		title = m.Title()
		if title == "" {
			title = strings.TrimSuffix(filepath.Base(path), ext)
		}
		artistName = m.Artist()
		if artistName == "" {
			artistName = "Unknown Artist"
		}
		albumName = m.Album()
		if albumName == "" {
			albumName = "Unknown Album"
		}
		if y := m.Year(); y > 0 {
			year = &y
		}
		if n, _ := m.Track(); n > 0 {
			trackNum = &n
		}
		if d, _ := m.Disc(); d > 0 {
			discNum = &d
		}
	}

	format := strings.TrimPrefix(ext, ".")

	// Upsert artist.
	artistID, err := s.store.UpsertArtist(ctx, artistName)
	if err != nil {
		return err
	}

	// Upsert album.
	albumID, err := s.store.UpsertAlbum(ctx, artistID, albumName, year)
	if err != nil {
		return err
	}

	// Insert track.
	return s.store.InsertTrack(ctx, Track{
		AlbumID:     &albumID,
		ArtistID:    artistID,
		Title:       title,
		TrackNumber: trackNum,
		DiscNumber:  discNum,
		FilePath:    path,
		Format:      &format,
		SizeBytes:   &size,
	})
}
