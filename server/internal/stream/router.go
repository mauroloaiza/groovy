package stream

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

var mimeTypes = map[string]string{
	"mp3":  "audio/mpeg",
	"flac": "audio/flac",
	"ogg":  "audio/ogg",
	"m4a":  "audio/mp4",
	"aac":  "audio/aac",
	"wav":  "audio/wav",
	"opus": "audio/ogg; codecs=opus",
	"wma":  "audio/x-ms-wma",
}

type handler struct {
	db *sql.DB
}

func Router(db *sql.DB) http.Handler {
	h := &handler{db: db}
	r := chi.NewRouter()
	r.Get("/{id}", h.stream)
	return r
}

// GET /api/stream/:id
// Supports HTTP Range requests — seekable playback.
func (h *handler) stream(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var filePath, format string
	err = h.db.QueryRowContext(r.Context(),
		`SELECT file_path, COALESCE(format, '') FROM tracks WHERE id = ?`, id,
	).Scan(&filePath, &format)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	if mime, ok := mimeTypes[format]; ok {
		w.Header().Set("Content-Type", mime)
	}
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "no-cache")

	// http.ServeContent handles Range, If-Range, Content-Length automatically.
	http.ServeContent(w, r, filePath, time.Time{}, f)
}
