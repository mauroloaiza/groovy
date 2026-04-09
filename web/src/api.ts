const BASE = '/api'

export interface Artist {
  id: number
  name: string
}

export interface Album {
  id: number
  artist_id: number
  name: string
  year?: number
}

export interface Track {
  id: number
  album_id?: number
  artist_id: number
  title: string
  track_number?: number
  disc_number?: number
  duration_sec?: number
  format?: string
  size_bytes?: number
}

async function get<T>(path: string): Promise<T> {
  const res = await fetch(BASE + path)
  if (!res.ok) throw new Error(`${res.status} ${path}`)
  return res.json()
}

export const api = {
  artists: () => get<Artist[]>('/library/artists'),
  artist: (id: number) => get<Artist>(`/library/artists/${id}`),
  albumsByArtist: (artistId: number) => get<Album[]>(`/library/artists/${artistId}/albums`),
  tracksByAlbum: (albumId: number) => get<Track[]>(`/library/albums/${albumId}/tracks`),
  tracksByArtist: (artistId: number) => get<Track[]>(`/library/artists/${artistId}/tracks`),
  search: (q: string) => get<Track[]>(`/library/search?q=${encodeURIComponent(q)}`),
  scan: (dir?: string) =>
    fetch(`${BASE}/library/scan${dir ? `?dir=${encodeURIComponent(dir)}` : ''}`, {
      method: 'POST',
    }).then((r) => r.json()),
}

export const streamUrl = (trackId: number) => `/api/stream/${trackId}`
