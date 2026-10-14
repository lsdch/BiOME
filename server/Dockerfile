FROM golang:1.26 AS builder

WORKDIR /src

# COPY go.mod go.sum ./
COPY . .


RUN go mod download

RUN CGO_ENABLED=1 \
    GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/server .


FROM debian:bookworm-slim

RUN useradd --system --uid 10001 --no-create-home app

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY config.yml /app/config.yml

RUN mkdir -p /data/imports \
    && chown -R app:app /data

USER app

EXPOSE 8080

ENTRYPOINT ["/app/server"]