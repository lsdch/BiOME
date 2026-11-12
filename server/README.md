# Biome Backend

## Docker image 

From the repository root:

```sh
docker build -f server/docker/Dockerfile -t biome-server:local ./server
```

The local `server/huma` module is required by `go.mod` and must be present in
the build context, including in CI. CGO is retained for H3; both the build and
runtime stages use Debian Bookworm. The image includes CA certificates,
time zone data and assets, but no deployment configuration files.

## Kubernetes configuration

Mount a `config.yml` file from a ConfigMap read-only into `/etc/biome`.
This file is required. Start from `server/config.yml` and complete the
installation-specific settings. The image supplies these values:

| Variable | Value |
| --- | --- |
| `BIOME_ENV` | `prod` |
| `GIN_MODE` | `release` |
| `BIOME_API_API_HOST` | `0.0.0.0` |
| `BIOME_API_API_PORT` | `8080` |
| `BIOME_RAW_FILE_STORAGE_ROOT` | `/data/raw_files` |
| `BIOME_BOOTSTRAP_COUNTRIES_COUNTRIES_JSON_CACHE_PATH` | `/data/cache/countries.json` |

`BIOME_*` environment variables take precedence over configuration files.
Their names follow the nested YAML keys, separated by `_`; repetitions such
as `API_API_PORT` are intentional. Keys may be absent from the YAML file.

Inject secrets through `env[].valueFrom.secretKeyRef`, particularly:

```yaml
env:
  - name: BIOME_DB_DSN
    valueFrom:
      secretKeyRef:
        name: biome-backend
        key: database-dsn
  - name: BIOME_AUTH_TOKENS_JWT_SECRET_KEY
    valueFrom:
      secretKeyRef:
        name: biome-backend
        key: jwt-secret-key
  - name: BIOME_AUTH_TOKENS_REFRESH_TOKEN_PEPPER
    valueFrom:
      secretKeyRef:
        name: biome-backend
        key: refresh-token-pepper
```

Also provide the public URL (`BIOME_APP_PUBLIC_BASE_URL`), administrator
contact (`BIOME_INSTANCE_ADMIN_EMAIL`), API metadata (`BIOME_API_VERSION`,
`BIOME_API_API_CONTACT_NAME`, `BIOME_API_API_CONTACT_EMAIL`), API base path
(`BIOME_API_API_BASE_PATH`, usually `/api/v1`) and `BIOME_GBIF_GBIF_USER_AGENT`,
either in the configuration file or in the Deployment environment variables.
`JWT_SECRET_KEY` requires at least 32 characters, and `REFRESH_TOKEN_PEPPER`
requires at least 8.

The optional files `/etc/biome/config.prod.yml` and then
`/etc/biome/config.local.yml` override the base configuration, before
environment variables are applied. Configuration changes require a Pod
restart. SMTP settings already persisted in the database take precedence
over configuration initialization values.

The process runs with UID/GID `10001:10001`. Provide writable volumes with
matching permissions (`fsGroup: 10001` if supported by the storage driver):

| Path | Usage |
| --- | --- |
| `/data/raw_files` | Imported source files; shared persistent storage if using multiple replicas |
| `/data/cache` | Country cache; `emptyDir` can be used |
| `/data/openapi` | Specifications generated at startup; `emptyDir` can be used |
| `/app/assets` | Images and editable application icon; persistent storage if needed |

Mounting an empty volume at `/app/assets` hides the files supplied by the
image: initialize its contents before use. `/app/data` points to `/data` for
legacy relative paths. The Kubernetes Service must target port 8080 unless
`BIOME_API_API_PORT` is overridden.

The default command is
`/app/server -config-dir /etc/biome -openapi /data/openapi`. The Deployment's
`args` can override these options. Local image checks do not replace testing
with PostgreSQL and the cluster's actual volumes.
