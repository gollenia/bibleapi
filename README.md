# Bible API

Small Go API for serving Bible content from one or more SQLite databases.

The service:

- loads every supported SQLite database in `./db` at startup
- exposes book metadata for each translation
- exposes a chapter/verse schema
- returns verses for single verses, chapters, and verse ranges

The HTTP server runs on port `8080`.

## Project Layout

- `main.go` defines the Fiber routes
- `storage.go` loads SQLite files from `./db` and builds in-memory caches
- `parser.go` parses verse requests such as `john/3:16-18`
- `db/` contains translation databases such as `slt.db`

## Requirements

- Go 1.22+
- at least one SQLite database in `./db`

The app automatically treats each database filename as a translation key. For example, `db/slt.db` becomes `translation=slt`.

## Run Locally

Install dependencies:

```bash
make setup
```

Start the API:

```bash
make run
```

Or run it directly:

```bash
go run .
```

Build a binary:

```bash
make build
```

## Run With Docker

Build and start:

```bash
docker compose up --build
```

The container mounts `./db` into `/app/db`, so local database files are available to the API.

## Translation Selection

Use the `translation` query parameter to select a database:

```text
GET /books?translation=slt
```

If `translation` is omitted or unknown, the API falls back to the first loaded database.

## API Endpoints

### `GET /books`

Returns the list of books for a translation.

Example:

```bash
curl "http://localhost:8080/books?translation=slt"
```

Response shape:

```json
{
  "translation": "slt",
  "books": [
    {
      "id": 43,
      "short_name": "john",
      "long_name": "John",
      "section": "gospels",
      "order": 43,
      "chapters": 21,
      "testament": "nt",
      "number": 4
    }
  ]
}
```

### `GET /schema`

Returns chapter-to-verse counts for each book, sorted by canonical order.

Example:

```bash
curl "http://localhost:8080/schema?translation=slt"
```

Response shape:

```json
{
  "translation": "slt",
  "schema": [
    {
      "id": "gen",
      "name": "Genesis",
      "order": 1,
      "verses": {
        "1": 31,
        "2": 25
      }
    }
  ]
}
```

### `GET /verses/*`

Returns verse text for a book reference.

Supported formats:

- `/verses/john/3` for a full chapter
- `/verses/john/3:16` for a single verse
- `/verses/john/3:16-18` for a range in one chapter
- `/verses/john/3:16-4:3` for a range across chapters

Examples:

```bash
curl "http://localhost:8080/verses/john/3?translation=slt"
```

```bash
curl "http://localhost:8080/verses/john/3:16?translation=slt"
```

```bash
curl "http://localhost:8080/verses/john/3:16-18?translation=slt"
```

```bash
curl "http://localhost:8080/verses/john/3:16-4:3?translation=slt"
```

Response shape:

```json
{
  "translation": "slt",
  "book": "John",
  "verses": [
    {
      "book_id": 43,
      "chapter": 3,
      "verse": 16,
      "text": "For God so loved the world...",
      "linebreak": 0
    }
  ]
}
```

## Error Behavior

- invalid verse path returns `400`
- unknown book returns `404`
- database/query failures return `500`

## Notes

- Book identifiers in verse URLs use the lowercase `short_name` from the `books` table.
- Supported database extensions are `.db`, `.sqlite`, and `.sqlite3`.
- The API will not start if no valid SQLite database is found in `./db`.
