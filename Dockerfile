# --- STAGE 1: Kompilieren im Go-Container ---
FROM golang:1.26-alpine AS builder

# Build-Tools installieren (wird für manche OS-Pakete benötigt)
RUN apk add --no-cache git

WORKDIR /app

# Dependencies cachen (beschleunigt spätere Builds)
COPY go.mod go.sum ./
RUN go mod download

# Den restlichen Code kopieren
COPY . .

# Das Binary absolut statisch und optimiert für Linux kompilieren
# (-ldflags "-w -s" entfernt Debug-Infos und spart massig Platz)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bibleapi .

# --- STAGE 2: Das minimale Laufzeit-Image ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Binary und Bibeldatenbanken ins Laufzeit-Image übernehmen
COPY --from=builder /app/bibleapi .
COPY db/*.SQLite3 ./db/

# Den Port für Fiber nach außen freigeben
EXPOSE 8080

# API starten
CMD ["./bibleapi"]
