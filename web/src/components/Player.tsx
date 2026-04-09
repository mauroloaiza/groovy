import { useEffect, useRef, useState } from 'react'
import {
  Play, Pause, SkipBack, SkipForward,
  Volume2, VolumeX,
} from 'lucide-react'
import { usePlayerStore, currentTrack } from '../store'
import { streamUrl } from '../api'
import { clsx } from 'clsx'

export default function Player() {
  const store = usePlayerStore()
  const track = usePlayerStore(currentTrack)
  const playing = usePlayerStore((s) => s.playing)
  const artistName = usePlayerStore((s) => s.artistName)
  const albumName = usePlayerStore((s) => s.albumName)

  const audioRef = useRef<HTMLAudioElement>(null)
  const [progress, setProgress] = useState(0)   // 0–1
  const [duration, setDuration] = useState(0)
  const [volume, setVolume] = useState(1)
  const [muted, setMuted] = useState(false)
  const [dragging, setDragging] = useState(false)

  // Sync src when track changes
  useEffect(() => {
    const audio = audioRef.current
    if (!audio || !track) return
    audio.src = streamUrl(track.id)
    audio.load()
    if (playing) audio.play().catch(console.error)
  }, [track?.id])

  // Sync play/pause state
  useEffect(() => {
    const audio = audioRef.current
    if (!audio || !track) return
    if (playing) audio.play().catch(console.error)
    else audio.pause()
  }, [playing])

  // Volume / mute
  useEffect(() => {
    if (audioRef.current) {
      audioRef.current.volume = volume
      audioRef.current.muted = muted
    }
  }, [volume, muted])

  function onTimeUpdate() {
    const audio = audioRef.current
    if (!audio || dragging) return
    setProgress(audio.duration ? audio.currentTime / audio.duration : 0)
    setDuration(audio.duration || 0)
  }

  function onEnded() {
    store.next()
  }

  function seek(e: React.ChangeEvent<HTMLInputElement>) {
    const val = parseFloat(e.target.value)
    setProgress(val)
    if (audioRef.current && audioRef.current.duration) {
      audioRef.current.currentTime = val * audioRef.current.duration
    }
  }

  function fmt(sec: number) {
    if (!isFinite(sec)) return '0:00'
    const m = Math.floor(sec / 60)
    const s = Math.floor(sec % 60)
    return `${m}:${s.toString().padStart(2, '0')}`
  }

  if (!track) return null

  return (
    <div className="fixed bottom-0 inset-x-0 bg-surface-1/95 backdrop-blur border-t border-white/5 px-4 py-3 z-50">
      <audio
        ref={audioRef}
        onTimeUpdate={onTimeUpdate}
        onLoadedMetadata={onTimeUpdate}
        onEnded={onEnded}
        onMouseDown={() => setDragging(true)}
        onMouseUp={() => setDragging(false)}
      />

      <div className="max-w-screen-xl mx-auto flex items-center gap-4">
        {/* Track info */}
        <div className="w-56 shrink-0 min-w-0">
          <p className="text-sm font-medium text-white truncate">{track.title}</p>
          <p className="text-xs text-zinc-500 truncate">
            {artistName}{albumName ? ` · ${albumName}` : ''}
          </p>
        </div>

        {/* Controls + seek */}
        <div className="flex-1 flex flex-col items-center gap-1.5">
          <div className="flex items-center gap-5">
            <button
              onClick={store.prev}
              className="text-zinc-400 hover:text-white transition-colors"
            >
              <SkipBack size={18} />
            </button>

            <button
              onClick={() => store.setPlaying(!playing)}
              className="w-9 h-9 rounded-full bg-white flex items-center justify-center hover:scale-105 transition-transform"
            >
              {playing
                ? <Pause size={16} fill="black" className="text-black" />
                : <Play size={16} fill="black" className="text-black ml-0.5" />
              }
            </button>

            <button
              onClick={store.next}
              className="text-zinc-400 hover:text-white transition-colors"
            >
              <SkipForward size={18} />
            </button>
          </div>

          {/* Seek bar */}
          <div className="w-full flex items-center gap-2">
            <span className="text-xs text-zinc-600 w-9 text-right shrink-0">
              {fmt(progress * duration)}
            </span>
            <input
              type="range"
              min={0} max={1} step={0.001}
              value={progress}
              onChange={seek}
              className={clsx('flex-1 h-1 accent-accent cursor-pointer', 'range-sm')}
            />
            <span className="text-xs text-zinc-600 w-9 shrink-0">{fmt(duration)}</span>
          </div>
        </div>

        {/* Volume */}
        <div className="w-32 shrink-0 flex items-center gap-2">
          <button
            onClick={() => setMuted(!muted)}
            className="text-zinc-500 hover:text-white transition-colors"
          >
            {muted || volume === 0
              ? <VolumeX size={16} />
              : <Volume2 size={16} />
            }
          </button>
          <input
            type="range"
            min={0} max={1} step={0.01}
            value={muted ? 0 : volume}
            onChange={(e) => { setVolume(parseFloat(e.target.value)); setMuted(false) }}
            className="flex-1 h-1 accent-accent cursor-pointer"
          />
        </div>
      </div>
    </div>
  )
}
