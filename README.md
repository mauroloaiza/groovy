# Groovy

Self-hosted music server with on-demand streaming, live radio, and AutoDJ. No dependencies on Navidrome or Subsonic.

## Features

- **Music Library** — Scan directories, read ID3/FLAC/OGG/M4A metadata, album art
- **On-demand streaming** — Direct or transcoded playback, seek support
- **Radio / AutoDJ** — Continuous broadcast stream, smart track selection, crossfade
- **Listener requests** — Request songs, vote, anti-repeat lock
- **Web UI** — React + Vite + Tailwind + shadcn
- **Android app** — Kotlin + Jetpack Compose + ExoPlayer

## Stack

| Layer | Tech |
|---|---|
| API / Backend | Go |
| Audio engine | Go + FFmpeg |
| Broadcast | Icecast2 |
| Database | PostgreSQL |
| Frontend | React + Vite + Tailwind + shadcn |
| Android | Kotlin + Jetpack Compose |
| Deploy | Docker Compose |

## REST API

```
GET  /health
POST /api/library/scan?dir=/path
GET  /api/library/artists
GET  /api/library/artists/:id
GET  /api/library/artists/:id/albums
GET  /api/library/artists/:id/tracks
GET  /api/library/albums/:id
GET  /api/library/albums/:id/tracks
GET  /api/library/tracks/:id
GET  /api/library/search?q=query
```

## Development

```bash
cp .env.example .env
docker compose -f docker-compose.dev.yml up
```

- API:     http://localhost:8080
- Web:     http://localhost:5173
- Icecast: http://localhost:8000

## License

MIT
