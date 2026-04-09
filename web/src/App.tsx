import { useState, useEffect } from 'react'
import { api, Artist, Album, Track } from './api'
import Sidebar from './components/Sidebar'
import AlbumGrid from './components/AlbumGrid'
import TrackList from './components/TrackList'
import Player from './components/Player'

type View = 'albums' | 'tracks'

export default function App() {
  const [artist, setArtist] = useState<Artist | null>(null)
  const [album, setAlbum] = useState<Album | null>(null)
  const [view, setView] = useState<View>('albums')
  const [tracks, setTracks] = useState<Track[]>([])

  async function selectArtist(a: Artist) {
    setArtist(a)
    setAlbum(null)
    setView('albums')
    setTracks([])
  }

  async function selectAlbum(al: Album) {
    setAlbum(al)
    setView('tracks')
    const t = await api.tracksByAlbum(al.id)
    setTracks(t ?? [])
  }

  async function showAllTracks() {
    if (!artist) return
    setAlbum(null)
    setView('tracks')
    const t = await api.tracksByArtist(artist.id)
    setTracks(t ?? [])
  }

  return (
    <div className="flex h-screen bg-surface text-white overflow-hidden">
      <Sidebar selectedId={artist?.id ?? null} onSelect={selectArtist} />

      <main className="flex-1 overflow-y-auto pb-24">
        {!artist && (
          <div className="flex items-center justify-center h-full text-zinc-600">
            <p>Select an artist to start</p>
          </div>
        )}

        {artist && view === 'albums' && (
          <AlbumGrid
            artist={artist}
            selectedId={album?.id ?? null}
            onSelect={selectAlbum}
            onAllTracks={showAllTracks}
          />
        )}

        {artist && view === 'tracks' && (
          <div>
            {/* Back to albums */}
            <div className="px-6 pt-5">
              <button
                onClick={() => setView('albums')}
                className="text-xs text-zinc-500 hover:text-white transition-colors mb-4"
              >
                ← {artist.name}
              </button>
            </div>
            <TrackList
              tracks={tracks}
              artist={artist}
              album={album ?? undefined}
              title={album ? album.name : `All tracks — ${artist.name}`}
            />
          </div>
        )}
      </main>

      <Player />
    </div>
  )
}
