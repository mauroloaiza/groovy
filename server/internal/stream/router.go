package stream

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	r.Head("/{id}", h.stream)
	return r
}

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
		http.Error(w, "file not found on disk", http.StatusNotFound)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		http.Error(w, "stat error", http.StatusInternalServerError)
		return
	}
	size := info.Size()

	mime := "audio/mpeg"
	if m, ok := mimeTypes[format]; ok {
		mime = m
	}

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Disposition", "inline")

	rangeHeader := r.Header.Get("Range")

	// No Range header — serve full file
	if rangeHeader == "" {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodHead {
			io.Copy(w, f)
		}
		return
	}

	// Parse Range: bytes=start-end  OR  bytes=start-
	start, end, ok := parseRange(rangeHeader, size)
	if !ok {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	length := end - start + 1
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
	w.WriteHeader(http.StatusPartialContent)

	if r.Method != http.MethodHead {
		f.Seek(start, io.SeekStart)
		io.CopyN(w, f, length)
	}
}

// parseRange handles:
//
//	bytes=0-1023       → start=0,  end=1023
//	bytes=1024-        → start=1024, end=size-1   ← Chrome always sends this
//	bytes=-512         → start=size-512, end=size-1
func parseRange(header string, size int64) (start, end int64, ok bool) {
	header = strings.TrimPrefix(header, "bytes=")
	parts := strings.SplitN(header, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}

	startStr, endStr := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	switch {
	case startStr == "" && endStr != "":
		// suffix range: bytes=-N → last N bytes
		n, err := strconv.ParseInt(endStr, 10, 64)
		if err != nil || n <= 0 {
			return 0, 0, false
		}
		start = size - n
		end = size - 1

	case startStr != "" && endStr == "":
		// open range: bytes=N- → from N to end
		s, err := strconv.ParseInt(startStr, 10, 64)
		if err != nil {
			return 0, 0, false
		}
		start = s
		end = size - 1

	default:
		// explicit range: bytes=N-M
		s, err1 := strconv.ParseInt(startStr, 10, 64)
		e, err2 := strconv.ParseInt(endStr, 10, 64)
		if err1 != nil || err2 != nil {
			return 0, 0, false
		}
		start, end = s, e
	}

	if end >= size {
		end = size - 1
	}
	if start < 0 || start > end {
		return 0, 0, false
	}

	return start, end, true
}
