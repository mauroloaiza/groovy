package library

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	store    *Store
	scanner  *Scanner
	musicDir string
}

func Router(db *sql.DB, musicDir string) http.Handler {
	store := NewStore(db)
	h := &handler{
		store:    store,
		scanner:  NewScanner(store),
		musicDir: musicDir,
	}

	r := chi.NewRouter()
	r.Post("/scan", h.scan)
	r.Get("/artists", h.listArtists)
	r.Get("/artists/{id}", h.getArtist)
	r.Get("/artists/{id}/albums", h.listAlbumsByArtist)
	r.Get("/artists/{id}/tracks", h.listTracksByArtist)
	r.Get("/albums/{id}", h.getAlbum)
	r.Get("/albums/{id}/tracks", h.listTracksByAlbum)
	r.Get("/tracks/{id}", h.getTrack)
	r.Get("/search", h.search)

	return r
}

func (h *handler) scan(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = h.musicDir
	}
	if dir == "" {
		dir = os.Getenv("MUSIC_DIR")
	}
	if dir == "" {
		http.Error(w, "music dir not configured — pass ?dir=/path/to/music", http.StatusBadRequest)
		return
	}
	result := h.scanner.Scan(r.Context(), dir)
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) listArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := h.store.ListArtists(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, artists)
}

func (h *handler) getArtist(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	artist, err := h.store.GetArtist(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if artist == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, artist)
}

func (h *handler) listAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	albums, err := h.store.ListAlbumsByArtist(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, albums)
}

func (h *handler) listTracksByArtist(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	tracks, err := h.store.ListTracksByArtist(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tracks)
}

func (h *handler) getAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	album, err := h.store.GetAlbum(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if album == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, album)
}

func (h *handler) listTracksByAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	tracks, err := h.store.ListTracksByAlbum(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tracks)
}

func (h *handler) getTrack(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	track, err := h.store.GetTrack(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if track == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, track)
}

func (h *handler) search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "q is required", http.StatusBadRequest)
		return
	}
	tracks, err := h.store.SearchTracks(r.Context(), q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tracks)
}

func parseID(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
