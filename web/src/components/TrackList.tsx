import { Play, Music } from 'lucide-react'
import { Track, Artist, Album } from '../api'
import { usePlayerStore, currentTrack } from '../store'
import { clsx } from 'clsx'

interface Props {
  tracks: Track[]
  artist: Artist
  album?: Album
  title: string
}

function fmt(sec?: number) {
  if (!sec) return '—'
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

export default function TrackList({ tracks, artist, album, title }: Props) {
  const store = usePlayerStore()
  const playing = usePlayerStore((s) => s.playing)
  const active = usePlayerStore(currentTrack)

  function play(index: number) {
    store.playQueue(tracks, index, artist.name, album?.name ?? '')
  }

  return (
    <div className="p-6">
      <div className="mb-5">
        <h2 className="text-xl font-bold text-white">{title}</h2>
        <p className="text-sm text-zinc-500 mt-0.5">{tracks.length} tracks</p>
      </div>

      {/* Play all button */}
      {tracks.length > 0 && (
        <button
          onClick={() => play(0)}
          className="flex items-center gap-2 mb-5 px-4 py-2 bg-accent hover:bg-accent-hover rounded-full text-sm font-medium text-white transition-colors"
        >
          <Play size={14} fill="currentColor" />
          Play all
        </button>
      )}

      <div className="space-y-0.5">
        {tracks.map((track, i) => {
          const isActive = active?.id === track.id
          return (
            <div
              key={track.id}
              onDoubleClick={() => play(i)}
              className={clsx(
                'group flex items-center gap-3 px-3 py-2 rounded-md cursor-default transition-colors',
                isActive ? 'bg-accent/15' : 'hover:bg-surface-3',
              )}
            >
              {/* Index / play indicator */}
              <div className="w-6 text-center shrink-0">
                {isActive && playing ? (
                  <Music size={13} className="text-accent mx-auto animate-pulse" />
                ) : (
                  <>
                    <span className={clsx('text-xs group-hover:hidden', isActive ? 'text-accent' : 'text-zinc-600')}>
                      {track.track_number ?? i + 1}
                    </span>
                    <Play
                      size={13}
                      className="hidden group-hover:block text-white mx-auto cursor-pointer"
                      fill="currentColor"
                      onClick={() => play(i)}
                    />
                  </>
                )}
              </div>

              {/* Title */}
              <span
                className={clsx(
                  'flex-1 text-sm truncate',
                  isActive ? 'text-accent font-medium' : 'text-zinc-300',
                )}
              >
                {track.title}
              </span>

              {/* Duration */}
              <span className="text-xs text-zinc-600 shrink-0">{fmt(track.duration_sec)}</span>
            </div>
          )
        })}
      </div>
    </div>
  )
}
