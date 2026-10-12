# Church Player Backend

REST API in Go that powers Church Player, an app used during church services to organize worship sets and play songs from a local drive (a USB stick or external disk).

It manages albums, songs, worship services, sections and the moments inside each section, and keeps their display order. Song files are read from a folder next to the executable, so the whole app can run straight from a portable drive.

## Stack

- Go 1.22, `gorilla/mux` for routing, `negroni` middleware, CORS via `gorilla/handlers`
- SQLite (`mattn/go-sqlite3`)
- Layered structure: `domain` (entities and DTOs), `usecase` (business rules), `sqlite` (repositories), `handler` (HTTP)
- Repository tests with an in-memory SQLite database and `testify`

## Project layout

```
main.go                         HTTP routes and app startup
internal/config                 environment configuration (SERVER_PORT)
internal/domain                 entities and DTOs
internal/infrastructure/usecase business rules
internal/infrastructure/sqlite  SQLite repositories and tests
internal/infrastructure/handler HTTP handlers
internal/infrastructure/fileoperator  song file access on the local drive
schema.sql                      database schema
```

## Running locally

Requires Go 1.22+ and a C compiler (needed by `go-sqlite3`).

```bash
go build -o church_player_backend .
sqlite3 church.db < schema.sql      # create the database next to the executable
./church_player_backend             # listens on :8080 (override with SERVER_PORT)
curl localhost:8080/ping
```

Cross-compiling for Windows (the drive is usually used on a Windows PC):

```bash
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o church_player_backend.exe .
```

## Tests

```bash
go test ./...
```

## API

Main resources, all JSON:

| Resource | Endpoints |
|---|---|
| Albums | `POST /albums`, `GET /albums`, `GET /albums/{id}`, `PATCH /albums`, `DELETE /albums/{id}` |
| Songs | `POST /songs/{album_id}`, `GET /songs`, `GET /songs/{id}`, `PATCH /songs`, `PATCH /songs/orders`, `DELETE /songs/{id}` |
| Moments | `POST /moments`, `GET /moments/{id}`, `GET /moments/{id}/all`, `PATCH /moments`, `PATCH /moments/orders`, `DELETE /moments/{id}` |
| Sections, worships | same pattern under `/sections` and `/worships` |
| Files | `GET /file-control/{path}` serves a song file from the drive |

A ready-to-use Postman collection is in `Church_Player.postman_collection.json`.

The frontend (Angular) lives in [church-player](https://github.com/ViniNepo/church-player).
