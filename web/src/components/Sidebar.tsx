import { useEffect, useState } from 'react'
import { Music2, Search, RefreshCw } from 'lucide-react'
import { api, Artist } from '../api'
import { clsx } from 'clsx'

interface Props {
  selectedId: number | null
  onSelect: (artist: Artist) => void
}

export default function Sidebar({ selectedId, onSelect }: Props) {
  const [artists, setArtists] = useState<Artist[]>([])
  const [query, setQuery] = useState('')
  const [scanning, setScanning] = useState(false)

  useEffect(() => {
    api.artists().then(setArtists).catch(console.error)
  }, [])

  const filtered = query
    ? artists.filter((a) => a.name.toLowerCase().includes(query.toLowerCase()))
    : artists

  async function handleScan() {
    setScanning(true)
    try {
      await api.scan()
      const updated = await api.artists()
      setArtists(updated)
    } finally {
      setScanning(false)
    }
  }

  return (
    <aside className="flex flex-col w-64 shrink-0 bg-surface-1 border-r border-white/5 h-full">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-4 border-b border-white/5">
        <div className="flex items-center gap-2 text-accent font-semibold text-lg">
          <Music2 size={20} />
          Groovy
        </div>
        <button
          onClick={handleScan}
          title="Rescan library"
          className="text-zinc-500 hover:text-white transition-colors"
        >
          <RefreshCw size={15} className={scanning ? 'animate-spin' : ''} />
        </button>
      </div>

      {/* Search */}
      <div className="px-3 py-2">
        <div className="flex items-center gap-2 bg-surface-3 rounded-md px-3 py-1.5">
          <Search size={13} className="text-zinc-500 shrink-0" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Filter artists…"
            className="bg-transparent text-sm text-white placeholder-zinc-600 outline-none w-full"
          />
        </div>
      </div>

      {/* Artist list */}
      <nav className="flex-1 overflow-y-auto py-1">
        {filtered.map((a) => (
          <button
            key={a.id}
            onClick={() => onSelect(a)}
            className={clsx(
              'w-full text-left px-4 py-2 text-sm truncate transition-colors',
              selectedId === a.id
                ? 'bg-accent/20 text-accent font-medium'
                : 'text-zinc-400 hover:text-white hover:bg-white/5',
            )}
          >
            {a.name}
          </button>
        ))}
      </nav>

      <div className="px-4 py-2 text-xs text-zinc-600 border-t border-white/5">
        {artists.length} artists
      </div>
    </aside>
  )
}
