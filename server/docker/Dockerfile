FROM golang:1.26-bookworm AS builder

WORKDIR /src

# Includes the local Huma module referenced by go.mod.
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/server .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/goose github.com/pressly/goose/v3/cmd/goose


FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --system --gid 10001 app \
    && useradd --system --uid 10001 --gid app --no-create-home app

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /out/goose /app/goose
COPY db/migrations/*.sql /app/migrations/
COPY --chown=10001:10001 assets/ /app/assets/

RUN mkdir -p /etc/biome /data/raw_files /data/cache /data/openapi \
    && ln -s /data /app/data \
    && chown -R app:app /data

ENV GIN_MODE=release \
    BIOME_ENV=prod \
    BIOME_API_API_HOST=0.0.0.0 \
    BIOME_API_API_PORT=8080 \
    BIOME_RAW_FILE_STORAGE_ROOT=/data/raw_files \
    BIOME_BOOTSTRAP_COUNTRIES_COUNTRIES_JSON_CACHE_PATH=/data/cache/countries.json

USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["/app/server"]
CMD ["-config-dir", "/etc/biome", "-openapi", "/data/openapi"]
