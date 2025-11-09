# 1) Build frontend
FROM oven/bun:1-alpine AS web
WORKDIR /app
COPY web/ ./web/
WORKDIR /app/web
RUN bun install
RUN bun run build

# 2) Build backend
FROM golang:1.25-alpine AS go
WORKDIR /src
COPY backend/ ./backend/
WORKDIR /src/backend/cmd/server
RUN go build -o /out/app

# 3) Final image
FROM alpine:3.20
WORKDIR /app
COPY --from=go /out/app /app/app
COPY --from=web /app/web/dist /app/ui
VOLUME ["/app/data"]
EXPOSE 8080 8787
CMD ["/app/app", "--db", "data/pastes.db", "--ui", "/app/ui", "--addr", ":8080", "--admin", ":8787"]
