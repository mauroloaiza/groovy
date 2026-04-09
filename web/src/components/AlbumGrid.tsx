import { useEffect, useState } from 'react'
import { Disc3 } from 'lucide-react'
import { api, Album, Artist } from '../api'
import { clsx } from 'clsx'

interface Props {
  artist: Artist
  selectedId: number | null
  onSelect: (album: Album) => void
  onAllTracks: () => void
}

export default function AlbumGrid({ artist, selectedId, onSelect, onAllTracks }: Props) {
  const [albums, setAlbums] = useState<Album[]>([])

  useEffect(() => {
    setAlbums([])
    api.albumsByArtist(artist.id).then(setAlbums).catch(console.error)
  }, [artist.id])

  return (
    <div className="p-6">
      <div className="flex items-end justify-between mb-5">
        <div>
          <h1 className="text-2xl font-bold text-white">{artist.name}</h1>
          <p className="text-sm text-zinc-500 mt-0.5">{albums.length} albums</p>
        </div>
        <button
          onClick={onAllTracks}
          className="text-xs text-zinc-400 hover:text-accent transition-colors px-3 py-1.5 border border-white/10 rounded-md hover:border-accent/40"
        >
          All tracks
        </button>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {albums.map((album) => (
          <button
            key={album.id}
            onClick={() => onSelect(album)}
            className={clsx(
              'group text-left rounded-lg p-3 transition-colors border',
              selectedId === album.id
                ? 'bg-accent/10 border-accent/40'
                : 'bg-surface-2 border-transparent hover:bg-surface-3 hover:border-white/10',
            )}
          >
            {/* Cover placeholder */}
            <div className="w-full aspect-square rounded-md bg-surface-3 flex items-center justify-center mb-3 group-hover:bg-surface overflow-hidden">
              <Disc3 size={36} className="text-zinc-700 group-hover:text-zinc-600" />
            </div>
            <p className="text-sm font-medium text-white truncate">{album.name}</p>
            {album.year && <p className="text-xs text-zinc-500 mt-0.5">{album.year}</p>}
          </button>
        ))}
      </div>
    </div>
  )
}
