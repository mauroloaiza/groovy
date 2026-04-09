# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project scaffold (Go API, React web, Docker Compose dev environment)
- PostgreSQL schema: artists, albums, tracks with indexes
- Library scanner: recursive directory walk, reads ID3/FLAC/OGG tags via dhowden/tag
- Library store: upsert artists/albums, insert tracks, queries by artist/album
- REST API `/api/library/*`: artists, albums, tracks, search, POST /scan
- Auto-migration on startup (embedded SQL)
- CORS middleware
