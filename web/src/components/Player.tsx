import { useEffect, useRef, useState } from 'react'
import { Play, Pause, SkipBack, SkipForward, Volume2, VolumeX } from 'lucide-react'
import { usePlayerStore, currentTrack } from '../store'
import { streamUrl } from '../api'

export default function Player() {
  const store      = usePlayerStore()
  const track      = usePlayerStore(currentTrack)
  const playing    = usePlayerStore((s) => s.playing)
  const artistName = usePlayerStore((s) => s.artistName)
  const albumName  = usePlayerStore((s) => s.albumName)

  const audioRef   = useRef<HTMLAudioElement>(null)
  const [progress, setProgress] = useState(0)
  const [duration, setDuration] = useState(0)
  const [volume, setVolume]     = useState(1)
  const [muted, setMuted]       = useState(false)

  // Single effect: react to track change
  useEffect(() => {
    const audio = audioRef.current
    if (!audio) return
    if (!track) { audio.pause(); audio.src = ''; return }

    const url = streamUrl(track.id)
    console.log('[Player] loading track', track.id, url)
    audio.src = url
    audio.volume = volume
    audio.muted  = muted

    audio.play()
      .then(() => console.log('[Player] play() OK'))
      .catch((e) => console.error('[Player] play() FAILED:', e))

  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [track?.id])

  // Single effect: react to play/pause toggle (without track change)
  useEffect(() => {
    const audio = audioRef.current
    if (!audio || !track) return
    console.log('[Player] playing changed →', playing)
    if (playing) {
      audio.play().catch((e) => console.error('[Player] resume failed:', e))
    } else {
      audio.pause()
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [playing])

  // Volume
  useEffect(() => {
    const audio = audioRef.current
    if (!audio) return
    audio.volume = volume
    audio.muted  = muted
  }, [volume, muted])

  function onTimeUpdate() {
    const a = audioRef.current
    if (!a) return
    setProgress(a.duration ? a.currentTime / a.duration : 0)
    setDuration(a.duration || 0)
  }

  function seek(e: React.ChangeEvent<HTMLInputElement>) {
    const val = parseFloat(e.target.value)
    setProgress(val)
    const a = audioRef.current
    if (a && isFinite(a.duration)) a.currentTime = val * a.duration
  }

  function fmt(sec: number) {
    if (!isFinite(sec) || isNaN(sec)) return '0:00'
    return `${Math.floor(sec / 60)}:${String(Math.floor(sec % 60)).padStart(2, '0')}`
  }

  return (
    <>
      <audio
        ref={audioRef}
        onTimeUpdate={onTimeUpdate}
        onLoadedMetadata={onTimeUpdate}
        onEnded={() => store.next()}
        onError={(e) => console.error('[Player] audio error', (e.target as HTMLAudioElement).error)}
      />

      {track && (
        <div className="fixed bottom-0 inset-x-0 bg-[#0f0f0f]/95 backdrop-blur border-t border-white/5 px-4 py-3 z-50">
          <div className="max-w-screen-xl mx-auto flex items-center gap-4">

            {/* Info */}
            <div className="w-56 shrink-0 min-w-0">
              <p className="text-sm font-medium text-white truncate">{track.title}</p>
              <p className="text-xs text-zinc-500 truncate">
                {artistName}{albumName ? ` · ${albumName}` : ''}
              </p>
            </div>

            {/* Controls */}
            <div className="flex-1 flex flex-col items-center gap-1.5">
              <div className="flex items-center gap-5">
                <button onClick={store.prev} className="text-zinc-400 hover:text-white transition-colors">
                  <SkipBack size={18} />
                </button>
                <button
                  onClick={() => store.setPlaying(!playing)}
                  className="w-9 h-9 rounded-full bg-white flex items-center justify-center hover:scale-105 transition-transform"
                >
                  {playing
                    ? <Pause size={16} fill="black" className="text-black" />
                    : <Play  size={16} fill="black" className="text-black ml-0.5" />}
                </button>
                <button onClick={store.next} className="text-zinc-400 hover:text-white transition-colors">
                  <SkipForward size={18} />
                </button>
              </div>

              {/* Seek */}
              <div className="w-full flex items-center gap-2">
                <span className="text-xs text-zinc-600 w-9 text-right shrink-0">{fmt(progress * duration)}</span>
                <input type="range" min={0} max={1} step={0.001} value={progress} onChange={seek} className="flex-1 cursor-pointer" />
                <span className="text-xs text-zinc-600 w-9 shrink-0">{fmt(duration)}</span>
              </div>
            </div>

            {/* Volume */}
            <div className="w-32 shrink-0 flex items-center gap-2">
              <button onClick={() => setMuted(!muted)} className="text-zinc-500 hover:text-white transition-colors">
                {muted || volume === 0 ? <VolumeX size={16} /> : <Volume2 size={16} />}
              </button>
              <input
                type="range" min={0} max={1} step={0.01}
                value={muted ? 0 : volume}
                onChange={(e) => { setVolume(parseFloat(e.target.value)); setMuted(false) }}
                className="flex-1 cursor-pointer"
              />
            </div>

          </div>
        </div>
      )}
    </>
  )
}
