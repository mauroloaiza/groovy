import { create } from 'zustand'
import { Track } from './api'

interface PlayerState {
  queue: Track[]
  currentIndex: number
  playing: boolean

  // Context labels (resolved by UI)
  artistName: string
  albumName: string

  playQueue: (tracks: Track[], index: number, artistName?: string, albumName?: string) => void
  playTrack: (track: Track, artistName?: string, albumName?: string) => void
  setPlaying: (v: boolean) => void
  next: () => void
  prev: () => void
  goTo: (index: number) => void
}

export const usePlayerStore = create<PlayerState>((set, get) => ({
  queue: [],
  currentIndex: 0,
  playing: false,
  artistName: '',
  albumName: '',

  playQueue(tracks, index, artistName = '', albumName = '') {
    set({ queue: tracks, currentIndex: index, playing: true, artistName, albumName })
  },

  playTrack(track, artistName = '', albumName = '') {
    set({ queue: [track], currentIndex: 0, playing: true, artistName, albumName })
  },

  setPlaying(v) {
    set({ playing: v })
  },

  next() {
    const { queue, currentIndex } = get()
    if (currentIndex < queue.length - 1)
      set({ currentIndex: currentIndex + 1, playing: true })
  },

  prev() {
    const { currentIndex } = get()
    if (currentIndex > 0)
      set({ currentIndex: currentIndex - 1, playing: true })
  },

  goTo(index) {
    set({ currentIndex: index, playing: true })
  },
}))

export const currentTrack = (s: PlayerState): Track | null =>
  s.queue[s.currentIndex] ?? null
